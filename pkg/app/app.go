package app

import (
	"embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed chartTemplate
var chartTemplate embed.FS

type options struct {
	ChartName  string
	Deployment bool
	Configmap  bool
	Service    bool
	Ingress    bool
	Volumes    bool
	// hpa bool
	// pdb bool
	// secret bool
	// sts bool
}

type App struct {
	chartPath string
	opts      options
}

func NewApp(chartName string, chartPath string, deployment bool, configmap bool, service bool, ingress bool, volumes bool) *App {
	return &App{
		chartPath: chartPath,
		opts: options{
			ChartName:  chartName,
			Deployment: deployment,
			Configmap:  configmap,
			Service:    service,
			Ingress:    ingress,
			Volumes:    volumes,
		},
	}
}

func (a *App) GenerateChart() error {
	// create directories
	err := os.MkdirAll(a.chartPath+string(os.PathSeparator)+"templates", 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(a.chartPath+string(os.PathSeparator)+"templates/tests", 0755)
	if err != nil {
		return err
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

	if a.opts.Deployment {
		err = createFileFromTemplate("chartTemplate/templates/deployment.yaml", a.chartPath+string(os.PathSeparator)+"templates/deployment.yaml", a.opts)
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
	return nil
}

func createFileFromTemplate(templatePath string, outputPath string, opts options) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// list every files of chartTemplate
	files, err := chartTemplate.ReadDir("chartTemplate")
	if err != nil {
		fmt.Println("Error reading template directory:", err)
		return err
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}

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
	// outputFile, err := os.Create(outputPath)
	// if err != nil {
	// 	return err
	// }
	// defer outputFile.Close()

	// copy file templatePath to outputFile from chartTemplate FS
	content, err := chartTemplate.ReadFile(templatePath)
	if err != nil {
		fmt.Println("Error opening template:", err)
		return err
	}
	err = os.WriteFile(outputPath, content, 0644)
	return err
}
