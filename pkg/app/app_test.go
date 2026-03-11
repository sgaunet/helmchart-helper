package app

import (
	"errors"
	"strings"
	"testing"

	charterrors "github.com/sgaunet/helmchart-helper/pkg/errors"
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
			name: "volumes file generation",
			opts: options{
				ChartName: "test-chart",
				Volumes:   true,
			},
			wantErr: false,
			expectedFiles: []string{
				"test-path/templates/pvc.yaml",
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

// newTestApp creates an App with mocks for error path testing.
func newTestApp(mockFS *mocks.MockFileSystem, mockTP *mocks.MockTemplateProcessor, opts options) *App {
	return &App{
		chartPath:         "test-path",
		opts:              opts,
		fs:                mockFS,
		templateProcessor: mockTP,
		pathManager:       mocks.NewMockPathManager(),
		chartTemplateFS:   GetChartTemplate(),
	}
}

func TestApp_createDirectoryStructure_errors(t *testing.T) {
	tests := []struct {
		name      string
		opts      options
		setupErr  func(*mocks.MockFileSystem, *mocks.MockTemplateProcessor)
		errType   charterrors.ErrorType
	}{
		{
			name: "MkdirAll fails for templates directory",
			opts: options{ChartName: "test-chart"},
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["MkdirAll:test-path/templates"] = errors.New("permission denied")
			},
			errType: charterrors.FileSystemError,
		},
		{
			name: "MkdirAll fails for tests directory with service enabled",
			opts: options{ChartName: "test-chart", Service: true},
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["MkdirAll:test-path/templates/tests"] = errors.New("permission denied")
			},
			errType: charterrors.FileSystemError,
		},
		{
			name: "create test-connection.yaml fails with service enabled",
			opts: options{ChartName: "test-chart", Service: true},
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["Create:test-path/templates/tests/test-connection.yaml"] = errors.New("disk full")
			},
			errType: charterrors.TemplateError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystem()
			mockTP := mocks.NewMockTemplateProcessor()
			tt.setupErr(mockFS, mockTP)

			app := newTestApp(mockFS, mockTP, tt.opts)
			err := app.createDirectoryStructure()

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var chartErr *charterrors.ChartError
			if !errors.As(err, &chartErr) {
				t.Fatalf("expected ChartError, got %T: %v", err, err)
			}
			if chartErr.Type != tt.errType {
				t.Errorf("expected error type %s, got %s", tt.errType, chartErr.Type)
			}
		})
	}
}

func TestApp_createFileFromTemplate_errors(t *testing.T) {
	tests := []struct {
		name     string
		setupErr func(*mocks.MockFileSystem, *mocks.MockTemplateProcessor)
		errType  charterrors.ErrorType
	}{
		{
			name: "file creation fails",
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["Create:test-path/output.yaml"] = errors.New("disk full")
			},
			errType: charterrors.FileSystemError,
		},
		{
			name: "template parsing fails",
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ParseFS:chartTemplate/Chart.yaml"] = errors.New("invalid template")
			},
			errType: charterrors.TemplateError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystem()
			mockTP := mocks.NewMockTemplateProcessor()
			tt.setupErr(mockFS, mockTP)

			app := newTestApp(mockFS, mockTP, options{ChartName: "test-chart"})
			err := app.createFileFromTemplate("chartTemplate/Chart.yaml", "test-path/output.yaml")

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var chartErr *charterrors.ChartError
			if !errors.As(err, &chartErr) {
				t.Fatalf("expected ChartError, got %T: %v", err, err)
			}
			if chartErr.Type != tt.errType {
				t.Errorf("expected error type %s, got %s", tt.errType, chartErr.Type)
			}
		})
	}
}

func TestApp_copyFileFromTemplate_errors(t *testing.T) {
	tests := []struct {
		name     string
		setupErr func(*mocks.MockFileSystem, *mocks.MockTemplateProcessor)
		errType  charterrors.ErrorType
	}{
		{
			name: "template read fails",
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ReadFile:chartTemplate/templates/helpers.tpl"] = errors.New("read error")
			},
			errType: charterrors.TemplateError,
		},
		{
			name: "file write fails",
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["WriteFile:test-path/output.tpl"] = errors.New("write error")
			},
			errType: charterrors.FileSystemError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystem()
			mockTP := mocks.NewMockTemplateProcessor()
			tt.setupErr(mockFS, mockTP)

			app := newTestApp(mockFS, mockTP, options{ChartName: "test-chart"})
			err := app.copyFileFromTemplate("chartTemplate/templates/helpers.tpl", "test-path/output.tpl")

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var chartErr *charterrors.ChartError
			if !errors.As(err, &chartErr) {
				t.Fatalf("expected ChartError, got %T: %v", err, err)
			}
			if chartErr.Type != tt.errType {
				t.Errorf("expected error type %s, got %s", tt.errType, chartErr.Type)
			}
		})
	}
}

func TestApp_appendToFile_errors(t *testing.T) {
	tests := []struct {
		name     string
		setupErr func(*mocks.MockFileSystem, *mocks.MockTemplateProcessor)
		errType  charterrors.ErrorType
		errOp    string
	}{
		{
			name: "template read fails",
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ReadFile:chartTemplate/templates/NOTES-DEFAULT.txt"] = errors.New("read error")
			},
			errType: charterrors.TemplateError,
			errOp:   "read-template",
		},
		{
			name: "file open fails",
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["OpenFile:test-path/NOTES.txt"] = errors.New("open error")
			},
			errType: charterrors.FileSystemError,
			errOp:   "append-file",
		},
		{
			name: "file write string fails",
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.WriteStringErrors["test-path/NOTES.txt"] = errors.New("write failed")
			},
			errType: charterrors.FileSystemError,
			errOp:   "append-file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystem()
			mockTP := mocks.NewMockTemplateProcessor()
			tt.setupErr(mockFS, mockTP)

			app := newTestApp(mockFS, mockTP, options{ChartName: "test-chart"})
			err := app.appendToFile("chartTemplate/templates/NOTES-DEFAULT.txt", "test-path/NOTES.txt")

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var chartErr *charterrors.ChartError
			if !errors.As(err, &chartErr) {
				t.Fatalf("expected ChartError, got %T: %v", err, err)
			}
			if chartErr.Type != tt.errType {
				t.Errorf("expected error type %s, got %s", tt.errType, chartErr.Type)
			}
			if chartErr.Operation != tt.errOp {
				t.Errorf("expected operation %s, got %s", tt.errOp, chartErr.Operation)
			}
		})
	}
}

func TestApp_generateBasicFiles_errors(t *testing.T) {
	tests := []struct {
		name     string
		setupErr func(*mocks.MockFileSystem, *mocks.MockTemplateProcessor)
	}{
		{
			name: "helpers.tpl copy fails",
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ReadFile:chartTemplate/templates/helpers.tpl"] = errors.New("read error")
			},
		},
		{
			name: "helmignore copy fails",
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ReadFile:chartTemplate/helmignore"] = errors.New("read error")
			},
		},
		{
			name: "Chart.yaml creation fails",
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["Create:test-path/Chart.yaml"] = errors.New("disk full")
			},
		},
		{
			name: "values.yaml creation fails",
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["Create:test-path/values.yaml"] = errors.New("disk full")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystem()
			mockTP := mocks.NewMockTemplateProcessor()
			tt.setupErr(mockFS, mockTP)

			app := newTestApp(mockFS, mockTP, options{ChartName: "test-chart"})
			err := app.generateBasicFiles()

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestApp_generateConditionalFiles_errors(t *testing.T) {
	tests := []struct {
		name     string
		opts     options
		setupErr func(*mocks.MockFileSystem, *mocks.MockTemplateProcessor)
	}{
		{
			name: "deployment template creation fails",
			opts: options{ChartName: "test-chart", Deployment: true},
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["Create:test-path/templates/deployment.yaml"] = errors.New("disk full")
			},
		},
		{
			name: "service template parsing fails",
			opts: options{ChartName: "test-chart", Service: true},
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ParseFS:chartTemplate/templates/service.yaml"] = errors.New("bad template")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystem()
			mockTP := mocks.NewMockTemplateProcessor()
			tt.setupErr(mockFS, mockTP)

			app := newTestApp(mockFS, mockTP, tt.opts)
			err := app.generateConditionalFiles()

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestApp_generateNotesFiles_errors(t *testing.T) {
	tests := []struct {
		name     string
		opts     options
		setupErr func(*mocks.MockFileSystem, *mocks.MockTemplateProcessor)
	}{
		{
			name: "NOTES-objects-created.txt creation fails",
			opts: options{ChartName: "test-chart"},
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["Create:test-path/templates/NOTES.txt"] = errors.New("disk full")
			},
		},
		{
			name: "NOTES-DEFAULT.txt append fails",
			opts: options{ChartName: "test-chart"},
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ReadFile:chartTemplate/templates/NOTES-DEFAULT.txt"] = errors.New("read error")
			},
		},
		{
			name: "NOTES-INGRESS.txt append fails with ingress enabled",
			opts: options{ChartName: "test-chart", Ingress: true},
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ReadFile:chartTemplate/templates/NOTES-INGRESS.txt"] = errors.New("read error")
			},
		},
		{
			name: "NOTES-SERVICE.txt append fails with service enabled",
			opts: options{ChartName: "test-chart", Service: true},
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ReadFile:chartTemplate/templates/NOTES-SERVICE.txt"] = errors.New("read error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystem()
			mockTP := mocks.NewMockTemplateProcessor()
			tt.setupErr(mockFS, mockTP)

			app := newTestApp(mockFS, mockTP, tt.opts)
			err := app.generateNotesFiles()

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestApp_replaceExampleInAllFiles_errors(t *testing.T) {
	tests := []struct {
		name     string
		setupErr func(*mocks.MockFileSystem)
		errOp    string
	}{
		{
			name: "Walk fails",
			setupErr: func(fs *mocks.MockFileSystem) {
				fs.Errors["Walk:test-path"] = errors.New("walk error")
			},
			errOp: "replace-placeholder",
		},
		{
			name: "ReadFile fails during walk",
			setupErr: func(fs *mocks.MockFileSystem) {
				fs.Files["test-path/Chart.yaml"] = []byte("name: example")
				fs.Errors["ReadFile:test-path/Chart.yaml"] = errors.New("read error")
			},
			errOp: "replace-placeholder",
		},
		{
			name: "WriteFile fails during walk",
			setupErr: func(fs *mocks.MockFileSystem) {
				fs.Files["test-path/Chart.yaml"] = []byte("name: example")
				fs.Errors["WriteFile:test-path/Chart.yaml"] = errors.New("write error")
			},
			errOp: "replace-placeholder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystem()
			tt.setupErr(mockFS)

			app := newTestApp(mockFS, mocks.NewMockTemplateProcessor(), options{ChartName: "test-chart"})
			err := app.replaceExampleInAllFiles("test-path")

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var chartErr *charterrors.ChartError
			if !errors.As(err, &chartErr) {
				t.Fatalf("expected ChartError, got %T: %v", err, err)
			}
			if chartErr.Operation != tt.errOp {
				t.Errorf("expected operation %s, got %s", tt.errOp, chartErr.Operation)
			}
		})
	}
}

func TestApp_GenerateChart_errors(t *testing.T) {
	tests := []struct {
		name        string
		opts        options
		setupErr    func(*mocks.MockFileSystem, *mocks.MockTemplateProcessor)
		errContains string
	}{
		{
			name: "fails at directory creation",
			opts: options{ChartName: "test-chart"},
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["MkdirAll:test-path/templates"] = errors.New("no space")
			},
			errContains: "create-directory",
		},
		{
			name: "fails at basic file generation",
			opts: options{ChartName: "test-chart"},
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ReadFile:chartTemplate/templates/helpers.tpl"] = errors.New("read error")
			},
			errContains: "read-template",
		},
		{
			name: "fails at conditional file generation",
			opts: options{ChartName: "test-chart", Deployment: true},
			setupErr: func(fs *mocks.MockFileSystem, _ *mocks.MockTemplateProcessor) {
				fs.Errors["Create:test-path/templates/deployment.yaml"] = errors.New("disk full")
			},
			errContains: "create-file",
		},
		{
			name: "fails at notes generation",
			opts: options{ChartName: "test-chart"},
			setupErr: func(_ *mocks.MockFileSystem, tp *mocks.MockTemplateProcessor) {
				tp.Errors["ParseFS:chartTemplate/templates/NOTES-objects-created.txt"] = errors.New("bad template")
			},
			errContains: "parse-template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystem()
			mockTP := mocks.NewMockTemplateProcessor()
			tt.setupErr(mockFS, mockTP)

			app := newTestApp(mockFS, mockTP, tt.opts)
			err := app.GenerateChart()

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("expected error containing %q, got: %v", tt.errContains, err)
			}
		})
	}
}