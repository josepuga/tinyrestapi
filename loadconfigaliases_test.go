package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigAliases(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.ini")

	jsonFile, err := filepath.Abs("app/data.json")
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
host=localhost
port=8001
json=%s

[Book Shop]
paths=sales
methods=GET
field aliases=id,code|hash,signature|badalias|noun,name
`, jsonFile)

	if err := os.WriteFile(configFile, []byte(config), 0644); err != nil {
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

	// Check valid aliases: alias -> internal JSON field.
	wantAliases := map[string]string{
		"id":   "code",
		"hash": "signature",
		"noun": "name",
	}

	for alias, internal := range wantAliases {
		got, exists := handler.jsonData.aliasesToInternal[alias]

		if !exists {
			t.Errorf(
				"alias %q not found in aliasesToInternal",
				alias,
			)
			continue
		}

		if got != internal {
			t.Errorf(
				"aliasesToInternal[%q] = %q, want %q",
				alias,
				got,
				internal,
			)
		}
	}

	// A malformed alias must be ignored.
	if _, exists := handler.jsonData.aliasesToInternal["badalias"]; exists {
		t.Error(
			`malformed alias "badalias" should not be loaded`,
		)
	}

	// Check the reverse mapping too.
	wantReverseAliases := map[string]string{
		"code":      "id",
		"signature": "hash",
		"name":      "noun",
	}

	for internal, alias := range wantReverseAliases {
		got, exists := handler.jsonData.internalToAliases[internal]

		if !exists {
			t.Errorf(
				"internal field %q not found in internalToAliases",
				internal,
			)
			continue
		}

		if got != alias {
			t.Errorf(
				"internalToAliases[%q] = %q, want %q",
				internal,
				got,
				alias,
			)
		}
	}
}