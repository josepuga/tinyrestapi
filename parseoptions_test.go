package main

import (
	"os"
	"testing"
)

func TestParseOptions(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantInteractive bool
		wantConfig      string
		wantErr         bool
	}{
		{
			name: "no arguments",
			args: []string{"tinyrestapi"},
		},
		{
			name:            "short interactive",
			args:            []string{"tinyrestapi", "-i"},
			wantInteractive: true,
		},
		{
			name:            "long interactive",
			args:            []string{"tinyrestapi", "-interactive"},
			wantInteractive: true,
		},
		{
			name:       "short config",
			args:       []string{"tinyrestapi", "-c", "test.ini"},
			wantConfig: "test.ini",
		},
		{
			name:       "long config",
			args:       []string{"tinyrestapi", "-config", "test.ini"},
			wantConfig: "test.ini",
		},
		{
			name:            "interactive and config",
			args:            []string{"tinyrestapi", "-i", "-c", "test.ini"},
			wantInteractive: true,
			wantConfig:      "test.ini",
		},
		{
			name:    "unknown option",
			args:    []string{"tinyrestapi", "-unknown"},
			wantErr: true,
		},
	}

	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			opt, err := parseOptions()

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"parseOptions() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}

			if opt.Interactive != tt.wantInteractive {
				t.Errorf(
					"Interactive = %v, want %v",
					opt.Interactive,
					tt.wantInteractive,
				)
			}

			if opt.ConfigFile != tt.wantConfig {
				t.Errorf(
					"ConfigFile = %q, want %q",
					opt.ConfigFile,
					tt.wantConfig,
				)
			}
		})
	}
}
