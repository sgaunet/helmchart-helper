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
	chartName  string
	deployment bool
	configmap  bool
	service    bool
	ingress    bool
	// hpa bool
	// pdb bool
	// secret bool
	// sts bool
	// volumes bool
}

type App struct {
	chartPath string
	opts      options
}

func NewApp(chartName string, chartPath string, deployment bool, configmap bool, service bool, ingress bool) *App {
	return &App{
		chartPath: chartPath,
		opts: options{
			chartName:  chartName,
			deployment: deployment,
			configmap:  configmap,
			service:    service,
			ingress:    ingress,
		},
	}
}

func (a *App) GenerateChart() error {
	err := createFileFromTemplate("chartTemplate/README.md", a.chartPath+string(os.PathSeparator)+"Chart.yaml", a.opts)
	if err != nil {
		return err
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
