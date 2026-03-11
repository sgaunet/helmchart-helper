// Package app provides the core functionality for generating Helm charts.
//
// Architecture:
//   - Uses dependency injection for filesystem, template processing, and path operations
//   - Embeds chart templates using Go's embed package
//   - Processes templates using text/template with conditional resource generation
//
// Main Components:
//   - App: Main application struct coordinating chart generation
//   - options: Configuration for which Kubernetes resources to generate
//   - chartTemplate: Embedded filesystem containing Helm chart templates
//
// Chart Generation Flow:
//  1. Create directory structure (chart root + templates/)
//  2. Generate basic files (Chart.yaml, values.yaml, _helpers.tpl, .helmignore)
//  3. Generate conditional resource files based on enabled options
//  4. Generate NOTES.txt with context-aware content
//  5. Replace "example" placeholder with actual chart name in all generated files
//
// Adding New Resource Types:
//  1. Add a bool field to the options struct
//  2. Add a Set<Resource> method on App
//  3. Add the template file to pkg/app/chartTemplate/templates/
//  4. Add the resource mapping in generateConditionalFiles()
//  5. Wire the new flag in pkg/cli/config.go and cmd/main.go
//
// Usage:
//
//	opts := app.NewApp("my-chart", "./output", fs, tp, pm, app.GetChartTemplate())
//	opts.SetDeployment(true)
//	opts.SetService(true)
//	err := opts.GenerateChart()
package app

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/sgaunet/helmchart-helper/pkg/errors"
	"github.com/sgaunet/helmchart-helper/pkg/interfaces"
)

//go:embed chartTemplate
var chartTemplate embed.FS

// GetChartTemplate returns the embedded chart template filesystem.
func GetChartTemplate() embed.FS {
	return chartTemplate
}

type options struct {
	ChartName      string
	Deployment     bool
	Cronjob        bool
	StatefulSet    bool
	DaemonSet      bool
	Configmap      bool
	Service        bool
	Ingress        bool
	Volumes        bool
	Hpa            bool
	ServiceAccount bool
	// pdb bool
	// secret bool
	// sts bool
}

// App manages Helm chart generation with configurable options.
type App struct {
	chartPath         string
	opts              options
	fs                interfaces.FileSystem
	templateProcessor interfaces.TemplateProcessor
	pathManager       interfaces.PathManager
	chartTemplateFS   embed.FS
}

// NewApp creates a new application instance for generating Helm charts.
func NewApp(chartName string, chartPath string, fs interfaces.FileSystem, templateProcessor interfaces.TemplateProcessor, pathManager interfaces.PathManager, chartTemplateFS embed.FS) *App {
	return &App{
		chartPath:         chartPath,
		fs:                fs,
		templateProcessor: templateProcessor,
		pathManager:       pathManager,
		chartTemplateFS:   chartTemplateFS,
		opts: options{
			ChartName: chartName,
		},
	}
}

// SetDeployment enables or disables Deployment resource generation.
func (a *App) SetDeployment(v bool) {
	a.opts.Deployment = v
}

// SetCronjob enables or disables CronJob resource generation.
func (a *App) SetCronjob(v bool) {
	a.opts.Cronjob = v
}

// SetDaemonSet enables or disables DaemonSet resource generation.
func (a *App) SetDaemonSet(v bool) {
	a.opts.DaemonSet = v
}

// SetConfigmap enables or disables ConfigMap resource generation.
func (a *App) SetConfigmap(v bool) {
	a.opts.Configmap = v
}

// SetService enables or disables Service resource generation.
func (a *App) SetService(v bool) {
	a.opts.Service = v
}

// SetIngress enables or disables Ingress resource generation.
func (a *App) SetIngress(v bool) {
	a.opts.Ingress = v
}

// SetVolumes enables or disables Volumes resource generation.
func (a *App) SetVolumes(v bool) {
	a.opts.Volumes = v
}

// SetHpa enables or disables HorizontalPodAutoscaler resource generation.
func (a *App) SetHpa(v bool) {
	a.opts.Hpa = v
}

// SetServiceAccount enables or disables ServiceAccount resource generation.
func (a *App) SetServiceAccount(v bool) {
	a.opts.ServiceAccount = v
}

// SetStatefulSet enables or disables StatefulSet resource generation.
func (a *App) SetStatefulSet(v bool) {
	a.opts.StatefulSet = v
}

// GenerateChart generates the complete Helm chart with all configured resources.
func (a *App) GenerateChart() error {
	if err := a.createDirectoryStructure(); err != nil {
		return err
	}
	
	if err := a.generateBasicFiles(); err != nil {
		return err
	}
	
	if err := a.generateConditionalFiles(); err != nil {
		return err
	}
	
	if err := a.generateNotesFiles(); err != nil {
		return err
	}
	
	if err := a.replaceTemplatePlaceholders(); err != nil {
		return err
	}
	
	return nil
}

func (a *App) createDirectoryStructure() error {
	// create directories
	templatesDir := a.chartPath + a.pathManager.Separator() + "templates"
	const dirPerm = 0755
	err := a.fs.MkdirAll(templatesDir, dirPerm)
	if err != nil {
		return errors.NewFileSystemError("create-directory", "failed to create templates directory", err).
			WithChart(a.opts.ChartName).
			WithFile(templatesDir)
	}
	
	if a.opts.Service {
		testsDir := a.chartPath + a.pathManager.Separator() + "templates/tests"
		err = a.fs.MkdirAll(testsDir, dirPerm)
		if err != nil {
			return errors.NewFileSystemError("create-directory", "failed to create tests directory", err).
				WithChart(a.opts.ChartName).
				WithFile(testsDir)
		}
		
		testFile := a.chartPath + a.pathManager.Separator() + "templates/tests/test-connection.yaml"
		err = a.createFileFromTemplate("chartTemplate/templates/tests/test-connection.yaml", testFile)
		if err != nil {
			return errors.WrapError(err, errors.TemplateError, "create-test-file", "failed to create test connection file").
				WithChart(a.opts.ChartName).
				WithFile(testFile)
		}
	}
	return nil
}

func (a *App) generateBasicFiles() error {
	// create files
	err := a.copyFileFromTemplate("chartTemplate/templates/helpers.tpl", a.chartPath+a.pathManager.Separator()+"templates/_helpers.tpl")
	if err != nil {
		return err
	}
	err = a.copyFileFromTemplate("chartTemplate/helmignore", a.chartPath+a.pathManager.Separator()+".helmignore")
	if err != nil {
		return err
	}
	err = a.createFileFromTemplate("chartTemplate/Chart.yaml", a.chartPath+a.pathManager.Separator()+"Chart.yaml")
	if err != nil {
		return err
	}
	err = a.createFileFromTemplate("chartTemplate/values.yaml", a.chartPath+a.pathManager.Separator()+"values.yaml")
	if err != nil {
		return err
	}
	return nil
}

func (a *App) generateConditionalFiles() error {
	type resourceMapping struct {
		enabled bool
		template string
		outputFile string
	}
	
	templatesPath := a.chartPath + a.pathManager.Separator() + "templates/"
	resources := []resourceMapping{
		{a.opts.Cronjob, "chartTemplate/templates/cronjob.yaml", templatesPath + "cronjob.yaml"},
		{a.opts.Deployment, "chartTemplate/templates/deployment.yaml", templatesPath + "deployment.yaml"},
		{a.opts.DaemonSet, "chartTemplate/templates/daemonset.yaml", templatesPath + "daemonset.yaml"},
		{a.opts.Service, "chartTemplate/templates/service.yaml", templatesPath + "service.yaml"},
		{a.opts.Ingress, "chartTemplate/templates/ingress.yaml", templatesPath + "ingress.yaml"},
		{a.opts.Configmap, "chartTemplate/templates/configmap.yaml", templatesPath + "configmap.yaml"},
		{a.opts.ServiceAccount, "chartTemplate/templates/serviceaccount.yaml", templatesPath + "serviceaccount.yaml"},
		{a.opts.StatefulSet, "chartTemplate/templates/statefulset.yaml", templatesPath + "statefulset.yaml"},
		{a.opts.Hpa, "chartTemplate/templates/hpa.yaml", templatesPath + "hpa.yaml"},
	}
	
	for _, resource := range resources {
		if resource.enabled {
			if err := a.createFileFromTemplate(resource.template, resource.outputFile); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func (a *App) generateNotesFiles() error {
	err := a.createFileFromTemplate("chartTemplate/templates/NOTES-objects-created.txt", a.chartPath+a.pathManager.Separator()+"templates/NOTES.txt")
	if err != nil {
		return err
	}
	err = a.appendToFile("chartTemplate/templates/NOTES-DEFAULT.txt", a.chartPath+a.pathManager.Separator()+"templates/NOTES.txt")
	if err != nil {
		return err
	}
	if a.opts.Ingress {
		err = a.appendToFile("chartTemplate/templates/NOTES-INGRESS.txt", a.chartPath+a.pathManager.Separator()+"templates/NOTES.txt")
		if err != nil {
			return err
		}
	}
	if a.opts.Service {
		err = a.appendToFile("chartTemplate/templates/NOTES-SERVICE.txt", a.chartPath+a.pathManager.Separator()+"templates/NOTES.txt")
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) replaceTemplatePlaceholders() error {
	// replace example with chart name
	return a.replaceExampleInAllFiles(a.chartPath)
}

func (a *App) createFileFromTemplate(templatePath string, outputPath string) error {
	outputFile, err := a.fs.Create(outputPath)
	if err != nil {
		return errors.NewFileSystemError("create-file", "failed to create output file", err).
			WithChart(a.opts.ChartName).
			WithFile(outputPath).
			WithContext("template", templatePath)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	tmpl, err := a.templateProcessor.ParseFS(a.chartTemplateFS, templatePath)
	if err != nil {
		return errors.NewTemplateError("parse-template", "failed to parse template", err).
			WithChart(a.opts.ChartName).
			WithFile(templatePath)
	}
	
	err = tmpl.Execute(outputFile, a.opts)
	if err != nil {
		return errors.NewTemplateError("execute-template", "failed to execute template", err).
			WithChart(a.opts.ChartName).
			WithFile(outputPath).
			WithContext("template", templatePath)
	}
	return nil
}

func (a *App) copyFileFromTemplate(templatePath string, outputPath string) error {
	// copy file templatePath to outputFile from chartTemplate FS
	content, err := a.templateProcessor.ReadFile(a.chartTemplateFS, templatePath)
	if err != nil {
		return fmt.Errorf("error opening template %s: %w", templatePath, err)
	}
	const filePerm = 0644
	if err = a.fs.WriteFile(outputPath, content, filePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (a *App) appendToFile(templatePath string, outputPath string) error {
	content, err := a.templateProcessor.ReadFile(a.chartTemplateFS, templatePath)
	if err != nil {
		return fmt.Errorf("error opening template %s: %w", templatePath, err)
	}

	const filePerm = 0644
	f, err := a.fs.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePerm)
	if err != nil {
		return fmt.Errorf("failed to open file for append: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err := f.WriteString(string(content)); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

// replaceExampleInAllFiles walks the generated chart directory and replaces all
// occurrences of the "example" placeholder with the actual chart name.
//
// The embedded chart templates use "example" as a placeholder name throughout
// (e.g., in Chart.yaml metadata, template helper names, NOTES.txt references).
// This function performs a post-generation pass to substitute the real chart name.
//
// It walks the directory tree using fs.Walk, skipping directories and processing
// only regular files. Each file is read, modified in-memory, and written back
// with a file permission of 0 (preserving existing permissions on overwrite).
func (a *App) replaceExampleInAllFiles(path string) error {
	err := a.fs.Walk(path, func(p string, info os.FileInfo, erR error) error {
		if erR != nil {
			return erR
		}
		if info.IsDir() {
			return nil
		}
		read, err := a.fs.ReadFile(p)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", p, err)
		}
		newContents := strings.ReplaceAll(string(read), "example", a.opts.ChartName)
		if err = a.fs.WriteFile(p, []byte(newContents), 0); err != nil {
			return fmt.Errorf("failed to write file %s: %w", p, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}
	return nil
}
