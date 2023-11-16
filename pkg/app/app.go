package app

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed chartTemplate
var chartTemplate embed.FS

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
	chartPath string
	opts      options
}

func NewApp(chartName string, chartPath string) *App {
	return &App{
		chartPath: chartPath,
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
	// create directories
	err := os.MkdirAll(a.chartPath+string(os.PathSeparator)+"templates", 0755)
	if err != nil {
		return err
	}
	if a.opts.Service {
		err = os.MkdirAll(a.chartPath+string(os.PathSeparator)+"templates/tests", 0755)
		if err != nil {
			return err
		}
		err = createFileFromTemplate("chartTemplate/templates/tests/test-connection.yaml", a.chartPath+string(os.PathSeparator)+"templates/tests/test-connection.yaml", a.opts)
		if err != nil {
			return err
		}
	}

	// create files
	err = copyFileFromTemplate("chartTemplate/templates/helpers.tpl", a.chartPath+string(os.PathSeparator)+"templates/_helpers.tpl")
	if err != nil {
		return err
	}
	err = copyFileFromTemplate("chartTemplate/helmignore", a.chartPath+string(os.PathSeparator)+".helmignore")
	if err != nil {
		return err
	}
	err = createFileFromTemplate("chartTemplate/Chart.yaml", a.chartPath+string(os.PathSeparator)+"Chart.yaml", a.opts)
	if err != nil {
		return err
	}
	err = createFileFromTemplate("chartTemplate/values.yaml", a.chartPath+string(os.PathSeparator)+"values.yaml", a.opts)
	if err != nil {
		return err
	}

	if a.opts.Cronjob {
		err = createFileFromTemplate("chartTemplate/templates/cronjob.yaml", a.chartPath+string(os.PathSeparator)+"templates/cronjob.yaml", a.opts)
		if err != nil {
			return err
		}
	}
	if a.opts.Deployment {
		err = createFileFromTemplate("chartTemplate/templates/deployment.yaml", a.chartPath+string(os.PathSeparator)+"templates/deployment.yaml", a.opts)
		if err != nil {
			return err
		}
	}
	if a.opts.DaemonSet {
		err = createFileFromTemplate("chartTemplate/templates/daemonset.yaml", a.chartPath+string(os.PathSeparator)+"templates/daemonset.yaml", a.opts)
		if err != nil {
			return err
		}
	}
	if a.opts.Service {
		err = createFileFromTemplate("chartTemplate/templates/service.yaml", a.chartPath+string(os.PathSeparator)+"templates/service.yaml", a.opts)
		if err != nil {
			return err
		}
	}
	if a.opts.Ingress {
		err = createFileFromTemplate("chartTemplate/templates/ingress.yaml", a.chartPath+string(os.PathSeparator)+"templates/ingress.yaml", a.opts)
		if err != nil {
			return err
		}
	}
	if a.opts.Configmap {
		err = createFileFromTemplate("chartTemplate/templates/configmap.yaml", a.chartPath+string(os.PathSeparator)+"templates/configmap.yaml", a.opts)
		if err != nil {
			return err
		}
	}
	if a.opts.ServiceAccount {
		err = createFileFromTemplate("chartTemplate/templates/serviceaccount.yaml", a.chartPath+string(os.PathSeparator)+"templates/serviceaccount.yaml", a.opts)
		if err != nil {
			return err
		}
	}
	if a.opts.StatefulSet {
		err = createFileFromTemplate("chartTemplate/templates/statefulset.yaml", a.chartPath+string(os.PathSeparator)+"templates/statefulset.yaml", a.opts)
		if err != nil {
			return err
		}
	}
	err = createFileFromTemplate("chartTemplate/templates/NOTES-objects-created.txt", a.chartPath+string(os.PathSeparator)+"templates/NOTES.txt", a.opts)
	if err != nil {
		return err
	}
	err = appendToFile("chartTemplate/templates/NOTES-DEFAULT.txt", a.chartPath+string(os.PathSeparator)+"templates/NOTES.txt")
	if err != nil {
		return err
	}
	if a.opts.Ingress {
		err = appendToFile("chartTemplate/templates/NOTES-INGRESS.txt", a.chartPath+string(os.PathSeparator)+"templates/NOTES.txt")
		if err != nil {
			return err
		}
	}
	if a.opts.Service {
		err = appendToFile("chartTemplate/templates/NOTES-SERVICE.txt", a.chartPath+string(os.PathSeparator)+"templates/NOTES.txt")
		if err != nil {
			return err
		}
	}
	// replace example with chart name
	err = a.replaceExampleInAllFiles(a.chartPath)
	return err
}

func createFileFromTemplate(templatePath string, outputPath string, opts options) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tmpl, err := template.ParseFS(chartTemplate, templatePath)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}
	err = tmpl.Execute(outputFile, opts)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}
	return nil
}

func copyFileFromTemplate(templatePath string, outputPath string) error {
	// copy file templatePath to outputFile from chartTemplate FS
	content, err := chartTemplate.ReadFile(templatePath)
	if err != nil {
		fmt.Println("Error opening template:", err)
		return err
	}
	err = os.WriteFile(outputPath, content, 0644)
	return err
}

func appendToFile(templatePath string, outputPath string) error {
	content, err := chartTemplate.ReadFile(templatePath)
	if err != nil {
		fmt.Println("Error opening template:", err)
		return err
	}

	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	err := filepath.Walk(path, func(p string, info os.FileInfo, erR error) error {
		// fmt.Println(p)
		// fmt.Println(info.Name())
		if erR != nil {
			return erR
		}
		if info.IsDir() {
			// return a.replaceExampleInAllFiles(p + string(os.PathSeparator) + info.Name())
			return nil
		}
		read, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		// fmt.Println(path)
		newContents := strings.Replace(string(read), "exemple", a.opts.ChartName, -1)
		err = os.WriteFile(p, []byte(newContents), 0)
		return err
	})
	return err
}
