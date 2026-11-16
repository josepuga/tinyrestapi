package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfigHandler(t *testing.T) {
	dir := t.TempDir()

	configFile := filepath.Join(dir, "config_test.ini")

	jsonFile, err := filepath.Abs("data_test.json")
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
host=localhost
port=8001
json=%s

[Book Shop]
paths=sales,ventas
methods=GET,POST,PUT
readonly=true
dry run=false
`, jsonFile)

	if err := os.WriteFile(
		configFile,
		[]byte(config),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	server := NewServer()

	if err := loadConfigFromFile(configFile, server); err != nil {
		t.Fatalf("loadConfigFromFile() error = %v", err)
	}

	if len(server.Handlers) != 1 {
		t.Fatalf(
			"len(server.Handlers) = %d, want 1",
			len(server.Handlers),
		)
	}

	handler := server.Handlers[0]

	if !handler.ReadOnlyMode {
		t.Error("handler.ReadOnlyMode = false, want true")
	}

	wantPaths := []string{
		"/sales",
		"/sales/:id",
		"/ventas",
		"/ventas/:id",
	}

	if !reflect.DeepEqual(handler.Paths, wantPaths) {
		t.Errorf(
			"handler.Paths = %v, want %v",
			handler.Paths,
			wantPaths,
		)
	}

	wantMethods := []string{
		"GET",
		"POST",
		"PUT",
	}

	for _, method := range wantMethods {
		if !handler.Methods[method] {
			t.Errorf(
				"handler.Methods[%q] = false, want true",
				method,
			)
		}
	}

	notWantedMethods := []string{
		"PATCH",
		"HEAD",
		"DELETE",
	}

	for _, method := range notWantedMethods {
		if handler.Methods[method] {
			t.Errorf(
				"handler.Methods[%q] = true, want false",
				method,
			)
		}
	}
}
