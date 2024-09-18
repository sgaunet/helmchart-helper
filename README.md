[![GitHub release](https://img.shields.io/github/release/sgaunet/helmchart-helper.svg)](https://github.com/sgaunet/helmchart-helper/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/helmchart-helper)](https://goreportcard.com/report/github.com/sgaunet/helmchart-helper)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/helmchart-helper/total)
[![Maintainability](https://api.codeclimate.com/v1/badges/3f62eea5e0ef5614c858/maintainability)](https://codeclimate.com/github/sgaunet/helmchart-helper/maintainability)
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

Download a release from the releases page.

or install it with go:

```bash
go install github.com/your-username/helmchart-helper@latest
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

## Issues and Bug Reports

We still encourage you to use our issue tracker for:

- 🐛 Reporting critical bugs
- 🔒 Reporting security vulnerabilities
- 🔍 Asking questions about the project

Please check existing issues before creating a new one to avoid duplicates.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
