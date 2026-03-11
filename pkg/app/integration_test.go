package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sgaunet/helmchart-helper/pkg/filesystem"
)

func TestGenerateChart_Integration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "helmchart-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name           string
		chartName      string
		options        map[string]bool
		expectedFiles  []string
		expectedDirs   []string
		fileChecks     map[string]func(string) error
	}{
		{
			name:      "basic chart generation",
			chartName: "my-test-chart",
			options:   map[string]bool{},
			expectedFiles: []string{
				"Chart.yaml",
				"values.yaml",
				".helmignore",
				"templates/_helpers.tpl",
				"templates/NOTES.txt",
			},
			expectedDirs: []string{
				"templates",
			},
			fileChecks: map[string]func(string) error{
				"Chart.yaml": func(content string) error {
					if !strings.Contains(content, "my-test-chart") {
						return &ValidationError{Field: "Chart.yaml", Message: "Chart name not found in Chart.yaml"}
					}
					return nil
				},
			},
		},
		{
			name:      "chart with deployment and service",
			chartName: "web-app",
			options: map[string]bool{
				"deployment": true,
				"service":    true,
			},
			expectedFiles: []string{
				"Chart.yaml",
				"values.yaml",
				".helmignore",
				"templates/_helpers.tpl",
				"templates/deployment.yaml",
				"templates/service.yaml",
				"templates/tests/test-connection.yaml",
				"templates/NOTES.txt",
			},
			expectedDirs: []string{
				"templates",
				"templates/tests",
			},
			fileChecks: map[string]func(string) error{
				"templates/deployment.yaml": func(content string) error {
					if !strings.Contains(content, "web-app") {
						return &ValidationError{Field: "deployment.yaml", Message: "Chart name not replaced in deployment template"}
					}
					return nil
				},
				"templates/service.yaml": func(content string) error {
					if !strings.Contains(content, "web-app") {
						return &ValidationError{Field: "service.yaml", Message: "Chart name not replaced in service template"}
					}
					return nil
				},
			},
		},
		{
			name:      "full featured chart",
			chartName: "full-app",
			options: map[string]bool{
				"deployment":     true,
				"service":        true,
				"ingress":        true,
				"configmap":      true,
				"serviceaccount": true,
				"hpa":            true,
				"volumes":        true,
			},
			expectedFiles: []string{
				"Chart.yaml",
				"values.yaml",
				".helmignore",
				"templates/_helpers.tpl",
				"templates/deployment.yaml",
				"templates/service.yaml",
				"templates/ingress.yaml",
				"templates/configmap.yaml",
				"templates/serviceaccount.yaml",
				"templates/hpa.yaml",
				"templates/pvc.yaml",
				"templates/tests/test-connection.yaml",
				"templates/NOTES.txt",
			},
			expectedDirs: []string{
				"templates",
				"templates/tests",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a subdirectory for this test
			testDir := filepath.Join(tempDir, tt.name)
			
			// Create dependencies with real filesystem
			fs := filesystem.NewOSFileSystem()
			templateProcessor := filesystem.NewDefaultTemplateProcessor()
			pathManager := filesystem.NewDefaultPathManager()
			
			// Create app
			app := NewApp(tt.chartName, testDir, fs, templateProcessor, pathManager, GetChartTemplate())
			
			// Set options
			if tt.options["deployment"] {
				app.SetDeployment(true)
			}
			if tt.options["service"] {
				app.SetService(true)
			}
			if tt.options["ingress"] {
				app.SetIngress(true)
			}
			if tt.options["configmap"] {
				app.SetConfigmap(true)
			}
			if tt.options["serviceaccount"] {
				app.SetServiceAccount(true)
			}
			if tt.options["hpa"] {
				app.SetHpa(true)
			}
			if tt.options["statefulset"] {
				app.SetStatefulSet(true)
			}
			if tt.options["daemonset"] {
				app.SetDaemonSet(true)
			}
			if tt.options["cronjob"] {
				app.SetCronjob(true)
			}
			if tt.options["volumes"] {
				app.SetVolumes(true)
			}

			// Generate chart
			err := app.GenerateChart()
			if err != nil {
				t.Errorf("GenerateChart() failed: %v", err)
				return
			}

			// Verify directories exist
			for _, expectedDir := range tt.expectedDirs {
				dirPath := filepath.Join(testDir, expectedDir)
				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					t.Errorf("Expected directory %s does not exist", expectedDir)
				}
			}

			// Verify files exist
			for _, expectedFile := range tt.expectedFiles {
				filePath := filepath.Join(testDir, expectedFile)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s does not exist", expectedFile)
				}
			}

			// Run file content checks
			for fileName, checkFunc := range tt.fileChecks {
				filePath := filepath.Join(testDir, fileName)
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read file %s: %v", fileName, err)
					continue
				}
				
				if err := checkFunc(string(content)); err != nil {
					t.Errorf("File validation failed for %s: %v", fileName, err)
				}
			}
		})
	}
}

func TestGenerateChart_ErrorHandling(t *testing.T) {
	// Use a read-only temp directory to trigger real filesystem errors
	tempDir, err := os.MkdirTemp("", "helmchart-error-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file where a directory is expected to block MkdirAll
	blockingFile := filepath.Join(tempDir, "blocked")
	if err := os.WriteFile(blockingFile, []byte("block"), 0644); err != nil {
		t.Fatalf("Failed to create blocking file: %v", err)
	}

	fs := filesystem.NewOSFileSystem()
	templateProcessor := filesystem.NewDefaultTemplateProcessor()
	pathManager := filesystem.NewDefaultPathManager()

	// Use the blocking file as the chart path; MkdirAll will fail
	// because "blocked/templates" can't be created under a regular file
	app := NewApp("test-chart", blockingFile, fs, templateProcessor, pathManager, GetChartTemplate())

	genErr := app.GenerateChart()
	if genErr == nil {
		t.Fatal("Expected error but got none")
	}

	if !strings.Contains(genErr.Error(), "create-directory") {
		t.Errorf("Expected error about directory creation, got: %v", genErr)
	}
}

// ValidationError represents a validation error during testing
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}