[![GitHub release](https://img.shields.io/github/release/sgaunet/helmchart-helper.svg)](https://github.com/sgaunet/helmchart-helper/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/helmchart-helper)](https://goreportcard.com/report/github.com/sgaunet/helmchart-helper)
![coverage](https://raw.githubusercontent.com/wiki/sgaunet/helmchart-helper/coverage-badge.svg)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/helmchart-helper/total)
[![License](https://img.shields.io/github/license/sgaunet/helmchart-helper.svg)](LICENSE)

# helmchart-helper

Helmchart-helper is a command-line tool for quickly generating basic Helm charts for your Kubernetes applications.

## Features

- Generation of the basic structure of a Helm chart
- Creation of essential files (Chart.yaml, values.yaml, templates, etc.)
- Simple customization via command-line options
- Generation of common Kubernetes manifests (Deployment, Service, Ingress)
- Validation of the generated chart

## Installation

### release

Download a release from the releases page and add the binary to your PATH.

### with go

```bash
go install github.com/your-username/helmchart-helper@latest
```

### homebrew

```bash
brew tap sgaunet/homebrew-tools
brew install sgaunet/tools/helmchart-helper
```

## Usage

```bash
Usage of helmchart-helper:
  -cj
        cronjob
  -cm
        configmap
  -deploy
        deployment
  -ds
        daemonse
  -help
        Print help
  -hpa
        hpa
  -ing
        ingress
  -n string
        Name of the chart
  -o string
        Path of the generated chart
  -pv
        volumes
  -sa
        serviceaccount
  -sts
        statefulset
  -svc
        service
  -version
        Print version
```

## 🕐 Project Status: Low Priority

This project is not under active development. While the project remains functional and available for use, please be aware of the following:

### What this means:
- **Response times will be longer** - Issues and pull requests may take weeks or months to be reviewed
- **Updates will be infrequent** - New features and non-critical bug fixes will be rare
- **Support is limited** - Questions and discussions may not receive timely responses

### We still welcome:
- 🐛 **Bug reports** - Critical issues will eventually be addressed
- 🔧 **Pull requests** - Well-tested contributions are appreciated
- 💡 **Feature requests** - Ideas will be considered for future development cycles
- 📖 **Documentation improvements** - Always helpful for the community

### Before contributing:
1. **Check existing issues** - Your concern may already be documented
2. **Be patient** - Responses may take considerable time
3. **Be self-sufficient** - Be prepared to fork and maintain your own version if needed
4. **Keep it simple** - Small, focused changes are more likely to be merged

### Alternative options:
If you need active support or rapid development:
- Look for actively maintained alternatives
- Reach out to discuss taking over maintenance

We appreciate your understanding and patience. This project remains important to us, but current priorities limit our ability to provide regular updates and support.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
