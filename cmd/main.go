package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sgaunet/helmchart-helper/pkg/app"
	// "github.com/bitfield/script"
)

var version = "dev"

func printVersion() {
	fmt.Printf("%s\n", version)
}

func main() {
	var (
		flagVersion bool
		flagHelp    bool
		chartName   string
		outputDir   string
	)

	flag.StringVar(&chartName, "n", "", "Name of the chart")
	flag.StringVar(&outputDir, "o", "", "Path of the generated chart")
	flag.BoolVar(&flagVersion, "version", false, "Print version")
	flag.BoolVar(&flagHelp, "help", false, "Print help")
	flag.Parse()

	if flagVersion {
		printVersion()
		os.Exit(0)
	}

	if flagHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// TODO: check if helm is installed
	if len(chartName) == 0 {
		fmt.Fprintf(os.Stderr, "Error: chart name is required\n")
		os.Exit(1)
	}

	if len(outputDir) == 0 {
		fmt.Fprintf(os.Stderr, "Error: chart path is required\n")
		os.Exit(1)
	}

	// TODO: check if chartPath exists (if not create it)

	// p := script.Exec("gum confirm \"Are you sure?\"").WithStderr(os.Stderr).WithStdout(os.Stdout)
	// p.Wait()
	// fmt.Println(p.ExitStatus())
	// chartName := "myChart"
	// chartPath := "tests/tmp/myChart"
	app := app.NewApp(chartName, outputDir, true, true, true, true)
	err := app.GenerateChart()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
