package cli

import (
	"strings"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			config: Config{
				ChartName: "test-chart",
				OutputDir: "/tmp/test",
			},
			wantErr: false,
		},
		{
			name: "missing chart name",
			config: Config{
				OutputDir: "/tmp/test",
			},
			wantErr:     true,
			errContains: "chart name is required",
		},
		{
			name: "missing output dir",
			config: Config{
				ChartName: "test-chart",
			},
			wantErr:     true,
			errContains: "chart path is required",
		},
		{
			name:        "empty config",
			config:      Config{},
			wantErr:     true,
			errContains: "chart name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Config.Validate() error message = %v, want to contain %v", err.Error(), tt.errContains)
			}
		})
	}
}

func TestConfig_HandleEarlyExit(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		version string
		want    bool
	}{
		{
			name: "version flag set",
			config: Config{
				Version: true,
			},
			version: "1.0.0",
			want:    true,
		},
		{
			name: "help flag set",
			config: Config{
				Help: true,
			},
			version: "1.0.0",
			want:    true,
		},
		{
			name: "no early exit flags",
			config: Config{
				ChartName: "test",
				OutputDir: "/tmp",
			},
			version: "1.0.0",
			want:    false,
		},
		{
			name: "both flags set - version takes precedence",
			config: Config{
				Version: true,
				Help:    true,
			},
			version: "1.0.0",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: We can't easily test the actual output without capturing stdout/stderr
			// but we can test the return value
			got := tt.config.HandleEarlyExit(tt.version)
			
			if got != tt.want {
				t.Errorf("Config.HandleEarlyExit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFlagsFromArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Config
	}{
		{
			name: "basic flags",
			args: []string{"-n", "test-chart", "-o", "/tmp/test"},
			expected: Config{
				ChartName: "test-chart",
				OutputDir: "/tmp/test",
			},
		},
		{
			name: "with deployment flag",
			args: []string{"-n", "test-chart", "-o", "/tmp/test", "-deploy"},
			expected: Config{
				ChartName:  "test-chart",
				OutputDir:  "/tmp/test",
				Deployment: true,
			},
		},
		{
			name: "multiple resource flags",
			args: []string{"-n", "test-chart", "-o", "/tmp/test", "-deploy", "-svc", "-ing"},
			expected: Config{
				ChartName:  "test-chart",
				OutputDir:  "/tmp/test",
				Deployment: true,
				Service:    true,
				Ingress:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse flags from args
			config, err := ParseFlagsFromArgs(tt.args)
			if err != nil {
				t.Errorf("ParseFlagsFromArgs() error = %v", err)
				return
			}

			// Check specific fields we care about
			if config.ChartName != tt.expected.ChartName {
				t.Errorf("ChartName = %v, want %v", config.ChartName, tt.expected.ChartName)
			}
			if config.OutputDir != tt.expected.OutputDir {
				t.Errorf("OutputDir = %v, want %v", config.OutputDir, tt.expected.OutputDir)
			}
			if config.Deployment != tt.expected.Deployment {
				t.Errorf("Deployment = %v, want %v", config.Deployment, tt.expected.Deployment)
			}
			if config.Service != tt.expected.Service {
				t.Errorf("Service = %v, want %v", config.Service, tt.expected.Service)
			}
			if config.Ingress != tt.expected.Ingress {
				t.Errorf("Ingress = %v, want %v", config.Ingress, tt.expected.Ingress)
			}
		})
	}
}