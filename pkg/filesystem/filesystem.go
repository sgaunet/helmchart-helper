// Package filesystem provides concrete implementations of filesystem interfaces.
package filesystem

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sgaunet/helmchart-helper/pkg/interfaces"
)

// OSFileSystem implements FileSystem interface using standard library.
type OSFileSystem struct{}

// NewOSFileSystem creates a new filesystem implementation using the OS.
func NewOSFileSystem() interfaces.FileSystem {
	return &OSFileSystem{}
}

// MkdirAll creates a directory named path, along with any necessary parents.
func (fs *OSFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	if err := os.MkdirAll(path, perm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// Create creates or truncates the named file.
func (fs *OSFileSystem) Create(name string) (interfaces.File, error) {
	f, err := os.Create(name) //nolint:gosec // G304: intentional file creation
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", name, err)
	}
	return f, nil
}

// WriteFile writes data to the named file, creating it if necessary.
func (fs *OSFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if err := os.WriteFile(name, data, perm); err != nil {
		return fmt.Errorf("failed to write file %s: %w", name, err)
	}
	return nil
}

// ReadFile reads the named file and returns the contents.
func (fs *OSFileSystem) ReadFile(name string) ([]byte, error) {
	data, err := os.ReadFile(name) //nolint:gosec // G304: intentional file reading
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", name, err)
	}
	return data, nil
}

// OpenFile opens the named file with specified flag.
func (fs *OSFileSystem) OpenFile(name string, flag int, perm fs.FileMode) (interfaces.File, error) {
	f, err := os.OpenFile(name, flag, perm) //nolint:gosec // G304: intentional file opening
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", name, err)
	}
	return f, nil
}

// Walk walks the file tree rooted at root, calling fn for each file or directory.
func (fs *OSFileSystem) Walk(root string, fn filepath.WalkFunc) error {
	if err := filepath.Walk(root, fn); err != nil {
		return fmt.Errorf("failed to walk directory %s: %w", root, err)
	}
	return nil
}

// DefaultTemplateProcessor implements TemplateProcessor interface.
type DefaultTemplateProcessor struct{}

// NewDefaultTemplateProcessor creates a new default template processor.
func NewDefaultTemplateProcessor() interfaces.TemplateProcessor { //nolint:ireturn // Returns interface by design
	return &DefaultTemplateProcessor{}
}

// ParseFS parses templates from the given embedded filesystem.
func (tp *DefaultTemplateProcessor) ParseFS(fs embed.FS, pattern string) (*template.Template, error) {
	tmpl, err := template.ParseFS(fs, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", pattern, err)
	}
	return tmpl, nil
}

// ReadFile reads a file from the embedded filesystem.
func (tp *DefaultTemplateProcessor) ReadFile(fs embed.FS, name string) ([]byte, error) {
	data, err := fs.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded file %s: %w", name, err)
	}
	return data, nil
}

// Execute applies a parsed template to the specified data object.
func (tp *DefaultTemplateProcessor) Execute(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	return buf.Bytes(), err
}

// DefaultPathManager implements PathManager interface.
type DefaultPathManager struct{}

// NewDefaultPathManager creates a new default path manager.
func NewDefaultPathManager() interfaces.PathManager { //nolint:ireturn // Returns interface by design
	return &DefaultPathManager{}
}

// Join joins path elements into a single path.
func (pm *DefaultPathManager) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Separator returns the OS-specific path separator.
func (pm *DefaultPathManager) Separator() string {
	return string(filepath.Separator)
}

// IsAbs reports whether the path is absolute.
func (pm *DefaultPathManager) IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

// Clean returns the shortest path name equivalent to path.
func (pm *DefaultPathManager) Clean(path string) string {
	return filepath.Clean(path)
}