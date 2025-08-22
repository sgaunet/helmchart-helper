package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/sgaunet/helmchart-helper/pkg/errors"
)

// Config holds all CLI configuration
type Config struct {
	ChartName      string
	OutputDir      string
	Deployment     bool
	Hpa            bool
	StatefulSet    bool
	DaemonSet      bool
	Cronjob        bool
	Configmap      bool
	Service        bool
	ServiceAccount bool
	Ingress        bool
	Volumes        bool
	Version        bool
	Help           bool
}

// ParseFlags parses command line flags and returns Config
func ParseFlags() (*Config, error) {
	return ParseFlagsFromArgs(os.Args[1:])
}

// ParseFlagsFromArgs parses flags from provided arguments (for testing)
func ParseFlagsFromArgs(args []string) (*Config, error) {
	config := &Config{}
	flagSet := flag.NewFlagSet("helmchart-helper", flag.ContinueOnError)
	
	flagSet.StringVar(&config.ChartName, "n", "", "Name of the chart")
	flagSet.StringVar(&config.OutputDir, "o", "", "Path of the generated chart")
	
	flagSet.BoolVar(&config.Hpa, "hpa", false, "hpa")
	flagSet.BoolVar(&config.StatefulSet, "sts", false, "statefulset")
	flagSet.BoolVar(&config.DaemonSet, "ds", false, "daemonse")
	flagSet.BoolVar(&config.Cronjob, "cj", false, "cronjob")
	flagSet.BoolVar(&config.Deployment, "deploy", false, "deployment")
	flagSet.BoolVar(&config.Configmap, "cm", false, "configmap")
	flagSet.BoolVar(&config.Ingress, "ing", false, "ingress")
	flagSet.BoolVar(&config.Volumes, "pv", false, "volumes")
	flagSet.BoolVar(&config.Service, "svc", false, "service")
	flagSet.BoolVar(&config.ServiceAccount, "sa", false, "serviceaccount")
	
	flagSet.BoolVar(&config.Version, "version", false, "Print version")
	flagSet.BoolVar(&config.Help, "help", false, "Print help")
	
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}
	
	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.ChartName == "" {
		return errors.NewValidationError("validate-config", "chart name is required").
			WithContext("flag", "-n")
	}
	
	if c.OutputDir == "" {
		return errors.NewValidationError("validate-config", "chart path is required").
			WithContext("flag", "-o")
	}
	
	return nil
}

// PrintVersion prints the version
func PrintVersion(version string) {
	fmt.Printf("%s\n", version)
}

// PrintHelp prints help information
func PrintHelp() {
	flag.PrintDefaults()
}

// HandleEarlyExit handles version and help flags that should exit early
func (c *Config) HandleEarlyExit(version string) bool {
	if c.Version {
		PrintVersion(version)
		return true
	}
	
	if c.Help {
		PrintHelp()
		return true
	}
	
	return false
}

// ExitWithError prints error and exits
func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// ExitSuccess exits with success code
func ExitSuccess() {
	os.Exit(0)
}