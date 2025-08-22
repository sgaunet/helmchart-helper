package filesystem

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sgaunet/helmchart-helper/pkg/interfaces"
)

// OSFileSystem implements FileSystem interface using standard library
type OSFileSystem struct{}

func NewOSFileSystem() interfaces.FileSystem {
	return &OSFileSystem{}
}

func (fs *OSFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *OSFileSystem) Create(name string) (interfaces.File, error) {
	return os.Create(name)
}

func (fs *OSFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (fs *OSFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fs *OSFileSystem) OpenFile(name string, flag int, perm fs.FileMode) (interfaces.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (fs *OSFileSystem) Walk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, fn)
}

// DefaultTemplateProcessor implements TemplateProcessor interface
type DefaultTemplateProcessor struct{}

func NewDefaultTemplateProcessor() interfaces.TemplateProcessor {
	return &DefaultTemplateProcessor{}
}

func (tp *DefaultTemplateProcessor) ParseFS(fs embed.FS, pattern string) (*template.Template, error) {
	return template.ParseFS(fs, pattern)
}

func (tp *DefaultTemplateProcessor) ReadFile(fs embed.FS, name string) ([]byte, error) {
	return fs.ReadFile(name)
}

func (tp *DefaultTemplateProcessor) Execute(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	return buf.Bytes(), err
}

// DefaultPathManager implements PathManager interface
type DefaultPathManager struct{}

func NewDefaultPathManager() interfaces.PathManager {
	return &DefaultPathManager{}
}

func (pm *DefaultPathManager) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (pm *DefaultPathManager) Separator() string {
	return string(filepath.Separator)
}

func (pm *DefaultPathManager) IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

func (pm *DefaultPathManager) Clean(path string) string {
	return filepath.Clean(path)
}