package main

import (
	"github.com/sgaunet/helmchart-helper/pkg/app"
	"github.com/sgaunet/helmchart-helper/pkg/cli"
	"github.com/sgaunet/helmchart-helper/pkg/filesystem"
)

var version = "dev"

func main() {
	// Parse CLI flags
	config, err := cli.ParseFlags()
	if err != nil {
		cli.ExitWithError(err)
	}

	// Handle early exit flags (version, help)
	if config.HandleEarlyExit(version) {
		cli.ExitSuccess()
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		cli.ExitWithError(err)
	}

	// Create dependencies
	fs := filesystem.NewOSFileSystem()
	templateProcessor := filesystem.NewDefaultTemplateProcessor()
	pathManager := filesystem.NewDefaultPathManager()
	
	// Create and configure app
	chartApp := app.NewApp(config.ChartName, config.OutputDir, fs, templateProcessor, pathManager, app.GetChartTemplate())
	chartApp.SetDeployment(config.Deployment)
	chartApp.SetHpa(config.Hpa)
	chartApp.SetStatefulSet(config.StatefulSet)
	chartApp.SetDaemonSet(config.DaemonSet)
	chartApp.SetCronjob(config.Cronjob)
	chartApp.SetConfigmap(config.Configmap)
	chartApp.SetIngress(config.Ingress)
	chartApp.SetVolumes(config.Volumes)
	chartApp.SetService(config.Service)
	chartApp.SetServiceAccount(config.ServiceAccount)

	// Generate chart
	if err := chartApp.GenerateChart(); err != nil {
		cli.ExitWithError(err)
	}
}
