package app

import (
	"testing"

	"github.com/sgaunet/helmchart-helper/pkg/mocks"
)

func TestApp_createDirectoryStructure(t *testing.T) {
	tests := []struct {
		name        string
		opts        options
		wantErr     bool
		expectedDirs map[string]bool
	}{
		{
			name: "basic directory structure",
			opts: options{
				ChartName: "test-chart",
				Service:   false,
			},
			wantErr: false,
			expectedDirs: map[string]bool{
				"test-path/templates": true,
			},
		},
		{
			name: "with service directory",
			opts: options{
				ChartName: "test-chart",
				Service:   true,
			},
			wantErr: false,
			expectedDirs: map[string]bool{
				"test-path/templates":       true,
				"test-path/templates/tests": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockFS := mocks.NewMockFileSystem()
			mockTemplateProcessor := mocks.NewMockTemplateProcessor()
			mockPathManager := mocks.NewMockPathManager()
			
			// Create app with mocks
			app := &App{
				chartPath:         "test-path",
				opts:              tt.opts,
				fs:                mockFS,
				templateProcessor: mockTemplateProcessor,
				pathManager:       mockPathManager,
				chartTemplateFS:   GetChartTemplate(),
			}

			// Execute
			err := app.createDirectoryStructure()

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("createDirectoryStructure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check directories were created
			for expectedDir := range tt.expectedDirs {
				if _, exists := mockFS.Directories[expectedDir]; !exists {
					t.Errorf("Expected directory %s was not created", expectedDir)
				}
			}
		})
	}
}

func TestApp_generateBasicFiles(t *testing.T) {
	tests := []struct {
		name             string
		opts             options
		wantErr          bool
		expectedFiles    []string
	}{
		{
			name: "basic files generation",
			opts: options{
				ChartName: "test-chart",
			},
			wantErr: false,
			expectedFiles: []string{
				"test-path/templates/_helpers.tpl",
				"test-path/.helmignore",
				"test-path/Chart.yaml",
				"test-path/values.yaml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockFS := mocks.NewMockFileSystem()
			mockTemplateProcessor := mocks.NewMockTemplateProcessor()
			mockPathManager := mocks.NewMockPathManager()
			
			// Create app with mocks
			app := &App{
				chartPath:         "test-path",
				opts:              tt.opts,
				fs:                mockFS,
				templateProcessor: mockTemplateProcessor,
				pathManager:       mockPathManager,
				chartTemplateFS:   GetChartTemplate(),
			}

			// Execute
			err := app.generateBasicFiles()

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("generateBasicFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check files were created
			for _, expectedFile := range tt.expectedFiles {
				if _, exists := mockFS.Files[expectedFile]; !exists {
					t.Errorf("Expected file %s was not created", expectedFile)
				}
			}
		})
	}
}

func TestApp_generateConditionalFiles(t *testing.T) {
	tests := []struct {
		name          string
		opts          options
		wantErr       bool
		expectedFiles []string
	}{
		{
			name: "deployment file generation",
			opts: options{
				ChartName:  "test-chart",
				Deployment: true,
			},
			wantErr: false,
			expectedFiles: []string{
				"test-path/templates/deployment.yaml",
			},
		},
		{
			name: "multiple conditional files",
			opts: options{
				ChartName:  "test-chart",
				Deployment: true,
				Service:    true,
				Ingress:    true,
			},
			wantErr: false,
			expectedFiles: []string{
				"test-path/templates/deployment.yaml",
				"test-path/templates/service.yaml",
				"test-path/templates/ingress.yaml",
			},
		},
		{
			name: "no conditional files",
			opts: options{
				ChartName: "test-chart",
			},
			wantErr:       false,
			expectedFiles: []string{}, // No files should be created
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockFS := mocks.NewMockFileSystem()
			mockTemplateProcessor := mocks.NewMockTemplateProcessor()
			mockPathManager := mocks.NewMockPathManager()
			
			// Create app with mocks
			app := &App{
				chartPath:         "test-path",
				opts:              tt.opts,
				fs:                mockFS,
				templateProcessor: mockTemplateProcessor,
				pathManager:       mockPathManager,
				chartTemplateFS:   GetChartTemplate(),
			}

			// Execute
			err := app.generateConditionalFiles()

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("generateConditionalFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check files were created
			for _, expectedFile := range tt.expectedFiles {
				if _, exists := mockFS.Files[expectedFile]; !exists {
					t.Errorf("Expected file %s was not created", expectedFile)
				}
			}
		})
	}
}

func TestApp_replaceTemplatePlaceholders(t *testing.T) {
	tests := []struct {
		name         string
		opts         options
		initialFiles map[string][]byte
		expectedContent string
		wantErr      bool
	}{
		{
			name: "replace example with chart name",
			opts: options{
				ChartName: "my-awesome-chart",
			},
			initialFiles: map[string][]byte{
				"test-path/Chart.yaml": []byte("name: example\nversion: 1.0.0"),
			},
			expectedContent: "name: my-awesome-chart\nversion: 1.0.0",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockFS := mocks.NewMockFileSystem()
			mockTemplateProcessor := mocks.NewMockTemplateProcessor()
			mockPathManager := mocks.NewMockPathManager()
			
			// Pre-populate files
			for path, content := range tt.initialFiles {
				mockFS.Files[path] = content
			}
			
			// Create app with mocks
			app := &App{
				chartPath:         "test-path",
				opts:              tt.opts,
				fs:                mockFS,
				templateProcessor: mockTemplateProcessor,
				pathManager:       mockPathManager,
				chartTemplateFS:   GetChartTemplate(),
			}

			// Execute
			err := app.replaceTemplatePlaceholders()

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceTemplatePlaceholders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check content was replaced
			if !tt.wantErr {
				for path := range tt.initialFiles {
					if content, exists := mockFS.Files[path]; exists {
						if string(content) != tt.expectedContent {
							t.Errorf("Expected content %s, got %s", tt.expectedContent, string(content))
						}
					}
				}
			}
		})
	}
}