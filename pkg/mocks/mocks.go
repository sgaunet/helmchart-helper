// Package mocks provides mock implementations for testing.
package mocks

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/sgaunet/helmchart-helper/pkg/interfaces"
)

var (
	// ErrFileNotFound is returned when a file is not found in the mock filesystem.
	ErrFileNotFound = errors.New("file not found")
)

// MockFileSystem implements FileSystem interface for testing.
type MockFileSystem struct {
	Files       map[string][]byte
	Directories map[string]bool
	Errors      map[string]error
}

// NewMockFileSystem creates a new mock file system for testing.
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files:       make(map[string][]byte),
		Directories: make(map[string]bool),
		Errors:      make(map[string]error),
	}
}

// MkdirAll simulates creating directories in the mock filesystem.
func (mfs *MockFileSystem) MkdirAll(path string, _ fs.FileMode) error {
	if err, exists := mfs.Errors["MkdirAll:"+path]; exists {
		return err
	}
	mfs.Directories[path] = true
	return nil
}

// Create simulates creating a file in the mock filesystem.
func (mfs *MockFileSystem) Create(name string) (interfaces.File, error) {
	if err, exists := mfs.Errors["Create:"+name]; exists {
		return nil, err
	}
	return &MockFile{name: name, fs: mfs}, nil
}

// WriteFile simulates writing to a file in the mock filesystem.
func (mfs *MockFileSystem) WriteFile(name string, data []byte, _ fs.FileMode) error {
	if err, exists := mfs.Errors["WriteFile:"+name]; exists {
		return err
	}
	mfs.Files[name] = data
	return nil
}

// ReadFile simulates reading a file from the mock filesystem.
func (mfs *MockFileSystem) ReadFile(name string) ([]byte, error) {
	if err, exists := mfs.Errors["ReadFile:"+name]; exists {
		return nil, err
	}
	if data, exists := mfs.Files[name]; exists {
		return data, nil
	}
	return nil, ErrFileNotFound
}

// OpenFile simulates opening a file in the mock filesystem.
func (mfs *MockFileSystem) OpenFile(name string, _ int, _ fs.FileMode) (interfaces.File, error) {
	if err, exists := mfs.Errors["OpenFile:"+name]; exists {
		return nil, err
	}
	return &MockFile{name: name, fs: mfs}, nil
}

// Walk simulates walking the directory tree in the mock filesystem.
func (mfs *MockFileSystem) Walk(root string, fn filepath.WalkFunc) error {
	if err, exists := mfs.Errors["Walk:"+root]; exists {
		return err
	}
	
	for path := range mfs.Files {
		if strings.HasPrefix(path, root) {
			info := &MockFileInfo{name: filepath.Base(path), isDir: false}
			if err := fn(path, info, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

// MockFile implements File interface for testing.
type MockFile struct {
	name string
	fs   *MockFileSystem
	buf  bytes.Buffer
}

// Write simulates writing bytes to the mock file.
func (mf *MockFile) Write(p []byte) (int, error) {
	n, err := mf.buf.Write(p)
	if err != nil {
		return n, fmt.Errorf("mock file write failed: %w", err)
	}
	return n, nil
}

// WriteString simulates writing a string to the mock file.
func (mf *MockFile) WriteString(s string) (int, error) {
	n, err := mf.buf.WriteString(s)
	if err != nil {
		return n, fmt.Errorf("mock file write string failed: %w", err)
	}
	return n, nil
}

// Close simulates closing the mock file.
func (mf *MockFile) Close() error {
	mf.fs.Files[mf.name] = mf.buf.Bytes()
	return nil
}

// MockFileInfo implements fs.FileInfo for testing.
type MockFileInfo struct {
	name  string
	isDir bool
}

// Name returns the name of the mock file.
func (mfi *MockFileInfo) Name() string       { return mfi.name }
// Size returns the size of the mock file.
func (mfi *MockFileInfo) Size() int64        { return 0 }
// Mode returns the file mode of the mock file.
func (mfi *MockFileInfo) Mode() fs.FileMode  { return 0 }
// ModTime returns the modification time of the mock file.
func (mfi *MockFileInfo) ModTime() time.Time { return time.Time{} }
// IsDir returns whether the mock file is a directory.
func (mfi *MockFileInfo) IsDir() bool        { return mfi.isDir }
// Sys returns the underlying system interface (always nil for mocks).
func (mfi *MockFileInfo) Sys() interface{}   { return nil }

// MockTemplateProcessor implements TemplateProcessor interface for testing.
type MockTemplateProcessor struct {
	Templates map[string]string
	Errors    map[string]error
}

// NewMockTemplateProcessor creates a new mock template processor for testing.
func NewMockTemplateProcessor() *MockTemplateProcessor {
	return &MockTemplateProcessor{
		Templates: make(map[string]string),
		Errors:    make(map[string]error),
	}
}

// ParseFS simulates parsing templates from an embedded filesystem.
func (mtp *MockTemplateProcessor) ParseFS(_ embed.FS, pattern string) (*template.Template, error) {
	if err, exists := mtp.Errors["ParseFS:"+pattern]; exists {
		return nil, err
	}
	if tmplStr, exists := mtp.Templates[pattern]; exists {
		tmpl, err := template.New("test").Parse(tmplStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse test template: %w", err)
		}
		return tmpl, nil
	}
	tmpl, err := template.New("test").Parse("{{.ChartName}}")
	if err != nil {
		return nil, fmt.Errorf("failed to parse default template: %w", err)
	}
	return tmpl, nil
}

// ReadFile simulates reading a file from an embedded filesystem.
func (mtp *MockTemplateProcessor) ReadFile(_ embed.FS, name string) ([]byte, error) {
	if err, exists := mtp.Errors["ReadFile:"+name]; exists {
		return nil, err
	}
	if content, exists := mtp.Templates[name]; exists {
		return []byte(content), nil
	}
	return []byte("mock content"), nil
}

// Execute simulates executing a template with data.
func (mtp *MockTemplateProcessor) Execute(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	return buf.Bytes(), err
}

// MockPathManager implements PathManager interface for testing.
type MockPathManager struct {
	JoinFunc      func(elem ...string) string
	SeparatorFunc func() string
	IsAbsFunc     func(path string) bool
	CleanFunc     func(path string) string
}

// NewMockPathManager creates a new mock path manager for testing.
func NewMockPathManager() *MockPathManager {
	return &MockPathManager{
		JoinFunc:      func(elem ...string) string { return strings.Join(elem, "/") },
		SeparatorFunc: func() string { return "/" },
		IsAbsFunc:     func(path string) bool { return strings.HasPrefix(path, "/") },
		CleanFunc:     func(path string) string { return path },
	}
}

// Join simulates joining path elements.
func (mpm *MockPathManager) Join(elem ...string) string {
	return mpm.JoinFunc(elem...)
}

// Separator returns the mock path separator.
func (mpm *MockPathManager) Separator() string {
	return mpm.SeparatorFunc()
}

// IsAbs simulates checking if a path is absolute.
func (mpm *MockPathManager) IsAbs(path string) bool {
	return mpm.IsAbsFunc(path)
}

// Clean simulates cleaning a path.
func (mpm *MockPathManager) Clean(path string) string {
	return mpm.CleanFunc(path)
}