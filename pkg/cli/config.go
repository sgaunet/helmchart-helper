// Package cli provides command line interface configuration for the Helm chart helper.
//
// It handles flag parsing, validation, and early-exit behaviors (--version, --help).
//
// Validation Constraints:
//   - Chart name (-n) must follow Helm naming conventions: start with a lowercase
//     letter, contain only lowercase letters, numbers, and hyphens, max 253 chars
//   - Output directory (-o) is required and must be non-empty
//   - All resource flags are optional and default to false
//
// Error Handling:
//   - Invalid flags return a wrapped error from flag.Parse
//   - Missing required fields return a ValidationError with flag context
//   - Early exit flags (--version, --help) are handled before validation
package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/sgaunet/helmchart-helper/pkg/errors"
)

// chartNameRegexp validates Helm chart names: must start with a lowercase letter,
// followed by lowercase letters, numbers, or hyphens. Cannot end with a hyphen.
var chartNameRegexp = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)

// maxChartNameLength is the maximum allowed chart name length (DNS subdomain limit).
const maxChartNameLength = 253

// Config holds all CLI configuration.
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

// ParseFlags parses command line flags and returns Config.
func ParseFlags() (*Config, error) {
	return ParseFlagsFromArgs(os.Args[1:])
}

// ParseFlagsFromArgs parses flags from provided arguments (for testing).
func ParseFlagsFromArgs(args []string) (*Config, error) {
	config := &Config{}
	flagSet := flag.NewFlagSet("helmchart-helper", flag.ContinueOnError)
	
	flagSet.StringVar(&config.ChartName, "n", "", "Name of the chart")
	flagSet.StringVar(&config.OutputDir, "o", "", "Path of the generated chart")
	
	flagSet.BoolVar(&config.Hpa, "hpa", false, "hpa")
	flagSet.BoolVar(&config.StatefulSet, "sts", false, "statefulset")
	flagSet.BoolVar(&config.DaemonSet, "ds", false, "daemonset")
	flagSet.BoolVar(&config.Cronjob, "cj", false, "cronjob")
	flagSet.BoolVar(&config.Deployment, "deploy", false, "deployment")
	flagSet.BoolVar(&config.Configmap, "cm", false, "configmap")
	flagSet.BoolVar(&config.Ingress, "ing", false, "ingress")
	flagSet.BoolVar(&config.Volumes, "pv", false, "volumes")
	flagSet.BoolVar(&config.Service, "svc", false, "service")
	flagSet.BoolVar(&config.ServiceAccount, "sa", false, "serviceaccount")
	
	flagSet.BoolVar(&config.Version, "version", false, "Print version")
	flagSet.BoolVar(&config.Help, "help", false, "Print help")
	
	if err := flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}
	
	return config, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if err := validateChartName(c.ChartName); err != nil {
		return err
	}

	if c.OutputDir == "" {
		return errors.NewValidationError("validate-config", "chart path is required").
			WithContext("flag", "-o")
	}

	return nil
}

// validateChartName validates a chart name against Helm naming conventions.
// Names must start with a lowercase letter, contain only lowercase letters,
// numbers, and hyphens, cannot end with a hyphen, and be at most 253 characters.
// Single-character names (a single lowercase letter) are also valid.
func validateChartName(name string) error {
	if name == "" {
		return errors.NewValidationError("validate-config", "chart name is required").
			WithContext("flag", "-n")
	}

	if len(name) > maxChartNameLength {
		return errors.NewValidationError("validate-config",
			fmt.Sprintf("chart name must be at most %d characters", maxChartNameLength)).
			WithContext("flag", "-n").
			WithContext("length", strconv.Itoa(len(name)))
	}

	// Single lowercase letter is valid
	if len(name) == 1 {
		if name[0] >= 'a' && name[0] <= 'z' {
			return nil
		}
		return errors.NewValidationError("validate-config",
			"chart name must start with a lowercase letter and contain only lowercase letters, numbers, and hyphens").
			WithContext("flag", "-n").
			WithContext("chart", name)
	}

	if !chartNameRegexp.MatchString(name) {
		return errors.NewValidationError("validate-config",
			"chart name must start with a lowercase letter and contain only lowercase letters, numbers, and hyphens").
			WithContext("flag", "-n").
			WithContext("chart", name)
	}

	return nil
}

// PrintVersion prints the version.
func PrintVersion(version string) {
	fmt.Printf("%s\n", version)
}

// PrintHelp prints help information.
func PrintHelp() {
	flag.PrintDefaults()
}

// HandleEarlyExit handles version and help flags that should exit early.
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

// ExitWithError prints error and exits.
func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// ExitSuccess exits with success code.
func ExitSuccess() {
	os.Exit(0)
}