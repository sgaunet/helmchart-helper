package mocks

import (
	"bytes"
	"embed"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/sgaunet/helmchart-helper/pkg/interfaces"
)

// MockFileSystem implements FileSystem interface for testing
type MockFileSystem struct {
	Files       map[string][]byte
	Directories map[string]bool
	Errors      map[string]error
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files:       make(map[string][]byte),
		Directories: make(map[string]bool),
		Errors:      make(map[string]error),
	}
}

func (mfs *MockFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	if err, exists := mfs.Errors["MkdirAll:"+path]; exists {
		return err
	}
	mfs.Directories[path] = true
	return nil
}

func (mfs *MockFileSystem) Create(name string) (interfaces.File, error) {
	if err, exists := mfs.Errors["Create:"+name]; exists {
		return nil, err
	}
	return &MockFile{name: name, fs: mfs}, nil
}

func (mfs *MockFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if err, exists := mfs.Errors["WriteFile:"+name]; exists {
		return err
	}
	mfs.Files[name] = data
	return nil
}

func (mfs *MockFileSystem) ReadFile(name string) ([]byte, error) {
	if err, exists := mfs.Errors["ReadFile:"+name]; exists {
		return nil, err
	}
	if data, exists := mfs.Files[name]; exists {
		return data, nil
	}
	return nil, errors.New("file not found")
}

func (mfs *MockFileSystem) OpenFile(name string, flag int, perm fs.FileMode) (interfaces.File, error) {
	if err, exists := mfs.Errors["OpenFile:"+name]; exists {
		return nil, err
	}
	return &MockFile{name: name, fs: mfs}, nil
}

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

// MockFile implements File interface for testing
type MockFile struct {
	name string
	fs   *MockFileSystem
	buf  bytes.Buffer
}

func (mf *MockFile) Write(p []byte) (int, error) {
	return mf.buf.Write(p)
}

func (mf *MockFile) WriteString(s string) (int, error) {
	return mf.buf.WriteString(s)
}

func (mf *MockFile) Close() error {
	mf.fs.Files[mf.name] = mf.buf.Bytes()
	return nil
}

// MockFileInfo implements fs.FileInfo for testing
type MockFileInfo struct {
	name  string
	isDir bool
}

func (mfi *MockFileInfo) Name() string       { return mfi.name }
func (mfi *MockFileInfo) Size() int64        { return 0 }
func (mfi *MockFileInfo) Mode() fs.FileMode  { return 0 }
func (mfi *MockFileInfo) ModTime() time.Time { return time.Time{} }
func (mfi *MockFileInfo) IsDir() bool        { return mfi.isDir }
func (mfi *MockFileInfo) Sys() interface{}   { return nil }

// MockTemplateProcessor implements TemplateProcessor interface for testing
type MockTemplateProcessor struct {
	Templates map[string]string
	Errors    map[string]error
}

func NewMockTemplateProcessor() *MockTemplateProcessor {
	return &MockTemplateProcessor{
		Templates: make(map[string]string),
		Errors:    make(map[string]error),
	}
}

func (mtp *MockTemplateProcessor) ParseFS(fs embed.FS, pattern string) (*template.Template, error) {
	if err, exists := mtp.Errors["ParseFS:"+pattern]; exists {
		return nil, err
	}
	if tmplStr, exists := mtp.Templates[pattern]; exists {
		return template.New("test").Parse(tmplStr)
	}
	return template.New("test").Parse("{{.ChartName}}")
}

func (mtp *MockTemplateProcessor) ReadFile(fs embed.FS, name string) ([]byte, error) {
	if err, exists := mtp.Errors["ReadFile:"+name]; exists {
		return nil, err
	}
	if content, exists := mtp.Templates[name]; exists {
		return []byte(content), nil
	}
	return []byte("mock content"), nil
}

func (mtp *MockTemplateProcessor) Execute(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	return buf.Bytes(), err
}

// MockPathManager implements PathManager interface for testing
type MockPathManager struct {
	JoinFunc      func(elem ...string) string
	SeparatorFunc func() string
	IsAbsFunc     func(path string) bool
	CleanFunc     func(path string) string
}

func NewMockPathManager() *MockPathManager {
	return &MockPathManager{
		JoinFunc:      func(elem ...string) string { return strings.Join(elem, "/") },
		SeparatorFunc: func() string { return "/" },
		IsAbsFunc:     func(path string) bool { return strings.HasPrefix(path, "/") },
		CleanFunc:     func(path string) string { return path },
	}
}

func (mpm *MockPathManager) Join(elem ...string) string {
	return mpm.JoinFunc(elem...)
}

func (mpm *MockPathManager) Separator() string {
	return mpm.SeparatorFunc()
}

func (mpm *MockPathManager) IsAbs(path string) bool {
	return mpm.IsAbsFunc(path)
}

func (mpm *MockPathManager) Clean(path string) string {
	return mpm.CleanFunc(path)
}