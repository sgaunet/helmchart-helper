package main

import (
	"fmt"
	"os"

	"github.com/sgaunet/helmchart-helper/pkg/app"
	// "github.com/bitfield/script"
)

func main() {
	// p := script.Exec("gum confirm \"Are you sure?\"").WithStderr(os.Stderr).WithStdout(os.Stdout)
	// p.Wait()
	// fmt.Println(p.ExitStatus())
	chartName := "myChart"
	chartPath := "tests/tmp/myChart"
	app := app.NewApp(chartName, chartPath, true, true, true, true)
	err := app.GenerateChart()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
