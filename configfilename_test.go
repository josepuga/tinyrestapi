package main

import (
	"testing"
)

func TestConfigFileName(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		env        string
		want       string
	}{
		{
			name:       "command line has priority",
			configFile: "cli.ini",
			env:        "env.ini",
			want:       "cli.ini",
		},
		{
			name: "environment variable",
			env:  "env.ini",
			want: "env.ini",
		},
		{
			name: "default config file",
			want: DefaultConfigFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TINYRESTAPICONFIG", tt.env)

			opt := Options{
				ConfigFile: tt.configFile,
			}

			got := configFileName(opt)

			if got != tt.want {
				t.Errorf(
					"configFileName() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}
