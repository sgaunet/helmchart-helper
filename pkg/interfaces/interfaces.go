// Package interfaces defines abstractions for filesystem, template processing, and path operations.
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
	Execute(tmpl *template.Template, data interface{}) ([]byte, error)
}

// PathManager abstracts path operations
type PathManager interface {
	Join(elem ...string) string
	Separator() string
	IsAbs(path string) bool
	Clean(path string) string
}