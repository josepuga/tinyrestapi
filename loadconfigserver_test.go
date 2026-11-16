package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigServer(t *testing.T) {
	dir := t.TempDir()

	jsonFile, err := filepath.Abs("data_test.json")
	if err != nil {
		t.Fatal(err)
	}

	configFile := filepath.Join(dir, "config_test.ini")

	config := fmt.Sprintf(`
host=localhost
port=8001
json=%s
`, jsonFile)

	if err := os.WriteFile(configFile, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	server := NewServer()

	err = loadConfigFromFile(configFile, server)
	if err != nil {
		t.Fatalf("loadConfigFromFile() error = %v", err)
	}

	if server.Port != 8001 {
		t.Errorf(
			"server.Port = %d, want %d",
			server.Port,
			8001,
		)
	}

	if server.Host != "localhost" {
		t.Errorf(
			"server.Host = %q, want %q",
			server.Host,
			"localhost",
		)
	}

	if server.Json != jsonFile {
		t.Errorf(
			"server.Host = %q, want %q",
			server.Json,
			jsonFile,
		)
	}
}
