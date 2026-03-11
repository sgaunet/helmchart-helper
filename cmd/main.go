// Package main is the entry point for the helmchart-helper CLI tool.
//
// It wires together the CLI flag parser (pkg/cli), production filesystem
// implementations (pkg/filesystem), and the chart generator (pkg/app).
//
// Execution flow:
//  1. Parse CLI flags → handle --version/--help → validate required flags
//  2. Create production dependencies (filesystem, template processor, path manager)
//  3. Configure the App with enabled resource types from CLI flags
//  4. Generate the Helm chart to the specified output directory
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
