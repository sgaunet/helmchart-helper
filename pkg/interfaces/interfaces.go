// Package interfaces defines abstractions for filesystem, template processing,
// and path operations used throughout the Helm chart helper.
//
// These interfaces enable dependency injection in the app package, allowing
// the core chart generation logic to be tested with mock implementations
// (see pkg/mocks) without touching the real filesystem.
//
// Main Interfaces:
//   - FileSystem: Abstracts directory creation, file read/write, and directory walking
//   - File: Abstracts individual file write and close operations
//   - TemplateProcessor: Abstracts Go template parsing and execution from embedded filesystems
//   - PathManager: Abstracts OS-specific path join operation
//
// Production implementations are in pkg/filesystem. Mock implementations are in pkg/mocks.
package interfaces

import (
	"embed"
	"io/fs"
	"path/filepath"
	"text/template"
)

// FileSystem abstracts file system operations for testing.
type FileSystem interface {
	MkdirAll(path string, perm fs.FileMode) error
	Create(name string) (File, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
	ReadFile(name string) ([]byte, error)
	OpenFile(name string, flag int, perm fs.FileMode) (File, error)
	Walk(root string, fn filepath.WalkFunc) error
}

// File abstracts file operations.
type File interface {
	Write(data []byte) (int, error)
	WriteString(s string) (int, error)
	Close() error
}

// TemplateProcessor abstracts template processing operations.
type TemplateProcessor interface {
	ParseFS(fs embed.FS, pattern string) (*template.Template, error)
	ReadFile(fs embed.FS, name string) ([]byte, error)
	Execute(tmpl *template.Template, data any) ([]byte, error)
}

// PathManager abstracts path operations.
type PathManager interface {
	Join(elem ...string) string
}