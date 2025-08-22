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

// GetChartTemplate returns the embedded chart template filesystem
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

type App struct {
	chartPath         string
	opts              options
	fs                interfaces.FileSystem
	templateProcessor interfaces.TemplateProcessor
	pathManager       interfaces.PathManager
	chartTemplateFS   embed.FS
}

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

func (a *App) SetDeployment(v bool) {
	a.opts.Deployment = v
}

func (a *App) SetCronjob(v bool) {
	a.opts.Cronjob = v
}

func (a *App) SetDaemonSet(v bool) {
	a.opts.DaemonSet = v
}

func (a *App) SetConfigmap(v bool) {
	a.opts.Configmap = v
}

func (a *App) SetService(v bool) {
	a.opts.Service = v
}

func (a *App) SetIngress(v bool) {
	a.opts.Ingress = v

}

func (a *App) SetVolumes(v bool) {
	a.opts.Volumes = v
}

func (a *App) SetHpa(v bool) {
	a.opts.Hpa = v
}

func (a *App) SetServiceAccount(v bool) {
	a.opts.ServiceAccount = v
}

func (a *App) SetStatefulSet(v bool) {
	a.opts.StatefulSet = v
}

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
	err := a.fs.MkdirAll(templatesDir, 0755)
	if err != nil {
		return errors.NewFileSystemError("create-directory", "failed to create templates directory", err).
			WithChart(a.opts.ChartName).
			WithFile(templatesDir)
	}
	
	if a.opts.Service {
		testsDir := a.chartPath + a.pathManager.Separator() + "templates/tests"
		err = a.fs.MkdirAll(testsDir, 0755)
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
	if a.opts.Cronjob {
		err := a.createFileFromTemplate("chartTemplate/templates/cronjob.yaml", a.chartPath+a.pathManager.Separator()+"templates/cronjob.yaml")
		if err != nil {
			return err
		}
	}
	if a.opts.Deployment {
		err := a.createFileFromTemplate("chartTemplate/templates/deployment.yaml", a.chartPath+a.pathManager.Separator()+"templates/deployment.yaml")
		if err != nil {
			return err
		}
	}
	if a.opts.DaemonSet {
		err := a.createFileFromTemplate("chartTemplate/templates/daemonset.yaml", a.chartPath+a.pathManager.Separator()+"templates/daemonset.yaml")
		if err != nil {
			return err
		}
	}
	if a.opts.Service {
		err := a.createFileFromTemplate("chartTemplate/templates/service.yaml", a.chartPath+a.pathManager.Separator()+"templates/service.yaml")
		if err != nil {
			return err
		}
	}
	if a.opts.Ingress {
		err := a.createFileFromTemplate("chartTemplate/templates/ingress.yaml", a.chartPath+a.pathManager.Separator()+"templates/ingress.yaml")
		if err != nil {
			return err
		}
	}
	if a.opts.Configmap {
		err := a.createFileFromTemplate("chartTemplate/templates/configmap.yaml", a.chartPath+a.pathManager.Separator()+"templates/configmap.yaml")
		if err != nil {
			return err
		}
	}
	if a.opts.ServiceAccount {
		err := a.createFileFromTemplate("chartTemplate/templates/serviceaccount.yaml", a.chartPath+a.pathManager.Separator()+"templates/serviceaccount.yaml")
		if err != nil {
			return err
		}
	}
	if a.opts.StatefulSet {
		err := a.createFileFromTemplate("chartTemplate/templates/statefulset.yaml", a.chartPath+a.pathManager.Separator()+"templates/statefulset.yaml")
		if err != nil {
			return err
		}
	}
	if a.opts.Hpa {
		err := a.createFileFromTemplate("chartTemplate/templates/hpa.yaml", a.chartPath+a.pathManager.Separator()+"templates/hpa.yaml")
		if err != nil {
			return err
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
	defer outputFile.Close()

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
		fmt.Println("Error opening template:", err)
		return err
	}
	err = a.fs.WriteFile(outputPath, content, 0644)
	return err
}

func (a *App) appendToFile(templatePath string, outputPath string) error {
	content, err := a.templateProcessor.ReadFile(a.chartTemplateFS, templatePath)
	if err != nil {
		fmt.Println("Error opening template:", err)
		return err
	}

	f, err := a.fs.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(string(content)); err != nil {
		return err
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
			return err
		}
		// fmt.Println(path)
		newContents := strings.Replace(string(read), "exemple", a.opts.ChartName, -1)
		err = a.fs.WriteFile(p, []byte(newContents), 0)
		return err
	})
	return err
}
