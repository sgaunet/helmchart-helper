// Package app provides the core functionality for generating Helm charts.
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

func (a *App) replaceExampleInAllFiles(path string) error {
	// list all files in chartPath
	err := a.fs.Walk(path, func(p string, info os.FileInfo, erR error) error {
		// fmt.Println(p)
		// fmt.Println(info.Name())
		if erR != nil {
			return erR
		}
		if info.IsDir() {
			// return a.replaceExampleInAllFiles(p + string(os.PathSeparator) + info.Name())
			return nil
		}
		read, err := a.fs.ReadFile(p)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", p, err)
		}
		// fmt.Println(path)
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
