package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfigHandlerWarnings(t *testing.T) {
	// Use an absolute path for the physical INI fixture so its location
	// does not depend on the current working directory.
	configFile, err := filepath.Abs("config_with_errors_test.ini")
	if err != nil {
		t.Fatal(err)
	}

	// The JSON filename inside the INI is relative (json=data.json).
	// For this fixture, data.json lives next to the INI, so temporarily
	// use the INI directory as the working directory.
	oldWorkingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(filepath.Dir(configFile)); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.Chdir(oldWorkingDir); err != nil {
			t.Errorf("unable to restore working directory: %v", err)
		}
	}()

	server := NewServer()

	// Invalid handler options are warnings, not fatal configuration errors.
	if err := loadConfigFromFile(configFile, server); err != nil {
		t.Fatalf("loadConfigFromFile() error = %v, want nil", err)
	}

	if len(server.Handlers) != 1 {
		t.Fatalf(
			"len(server.Handlers) = %d, want 1",
			len(server.Handlers),
		)
	}

	handler := server.Handlers[0]

	// Valid paths must still be loaded.
	wantPaths := []string{
		"/books",
		"/books/:id",
		"/libros",
		"/libros/:id",
	}

	if !reflect.DeepEqual(handler.Paths, wantPaths) {
		t.Errorf(
			"handler.Paths = %v, want %v",
			handler.Paths,
			wantPaths,
		)
	}

	// Valid methods must remain enabled.
	if !handler.Methods["GET"] {
		t.Error(`handler.Methods["GET"] = false, want true`)
	}

	if !handler.Methods["POST"] {
		t.Error(`handler.Methods["POST"] = false, want true`)
	}

	// Unknown methods must be ignored.
	if _, exists := handler.Methods["INVALID"]; exists {
		t.Error(`handler.Methods["INVALID"] exists, want unknown method ignored`)
	}

	// Valid aliases must still be loaded.
	if got := handler.jsonData.aliasesToInternal["id"]; got != "code" {
		t.Errorf(
			`aliasesToInternal["id"] = %q, want %q`,
			got,
			"code",
		)
	}

	if got := handler.jsonData.aliasesToInternal["name"]; got != "noun" {
		t.Errorf(
			`aliasesToInternal["name"] = %q, want %q`,
			got,
			"noun",
		)
	}

	// Malformed aliases must be ignored.
	if _, exists := handler.jsonData.aliasesToInternal["brokenalias"]; exists {
		t.Error(`malformed alias "brokenalias" should not be loaded`)
	}
}
