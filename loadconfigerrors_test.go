package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigErrors(t *testing.T) {
	tests := []struct {
		name   string
		config string
        createFile bool
	}{
		{
			name: "invalid port zero",
			config: `
host=localhost
port=0
json=data.json
`,
		},
		{
			name: "invalid port too high",
			config: `
host=localhost
port=70000
json=data.json
`,
		},
		{
			name: "missing json file",
			config: `
host=localhost
port=8001
json=missing.json
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			configFile := filepath.Join(dir, "config.ini")

			server := NewServer()

			if tt.createFile {
				// Make data.json resolvable for cases where the JSON file
				// itself is not the error being tested.
				jsonFile, err := filepath.Abs("app/data.json")
				if err != nil {
					t.Fatal(err)
				}

				config := fmt.Sprintf(tt.config, jsonFile)

				if err := os.WriteFile(
					configFile,
					[]byte(config),
					0644,
				); err != nil {
					t.Fatal(err)
				}
			}

			err := loadConfigFromFile(configFile, server)

			if err == nil {
				t.Fatal("loadConfigFromFile() expected error")
			}
		})
	}
}
