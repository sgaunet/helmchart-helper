// Package errors provides structured error handling for the Helm chart helper.
package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of error.
type ErrorType string

const (
	// ValidationError represents validation failures.
	ValidationError ErrorType = "validation"
	// FileSystemError represents file system operation failures.
	FileSystemError ErrorType = "filesystem"
	// TemplateError represents template processing failures.
	TemplateError ErrorType = "template"
	// ConfigurationError represents configuration-related errors.
	ConfigurationError ErrorType = "configuration"
)

// ChartError represents a structured error with context.
type ChartError struct {
	Type       ErrorType
	Operation  string
	Message    string
	Underlying error
	Context    map[string]string
}

// NewValidationError creates a new validation error.
func NewValidationError(operation, message string) *ChartError {
	return &ChartError{
		Type:      ValidationError,
		Operation: operation,
		Message:   message,
		Context:   make(map[string]string),
	}
}

// NewFileSystemError creates a new file system error.
func NewFileSystemError(operation, message string, underlying error) *ChartError {
	return &ChartError{
		Type:       FileSystemError,
		Operation:  operation,
		Message:    message,
		Underlying: underlying,
		Context:    make(map[string]string),
	}
}

// NewTemplateError creates a new template error.
func NewTemplateError(operation, message string, underlying error) *ChartError {
	return &ChartError{
		Type:       TemplateError,
		Operation:  operation,
		Message:    message,
		Underlying: underlying,
		Context:    make(map[string]string),
	}
}

// NewConfigurationError creates a new configuration error.
func NewConfigurationError(operation, message string) *ChartError {
	return &ChartError{
		Type:      ConfigurationError,
		Operation: operation,
		Message:   message,
		Context:   make(map[string]string),
	}
}

// Error implements the error interface.
func (e *ChartError) Error() string {
	var parts []string
	
	if e.Operation != "" {
		parts = append(parts, "operation: "+e.Operation)
	}
	
	if e.Type != "" {
		parts = append(parts, fmt.Sprintf("type: %s", e.Type))
	}
	
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	
	if e.Underlying != nil {
		parts = append(parts, fmt.Sprintf("caused by: %v", e.Underlying))
	}
	
	if len(e.Context) > 0 {
		var contextParts []string
		for k, v := range e.Context {
			contextParts = append(contextParts, fmt.Sprintf("%s=%s", k, v))
		}
		parts = append(parts, "context: "+strings.Join(contextParts, ", "))
	}
	
	return strings.Join(parts, "; ")
}

// Unwrap returns the underlying error for error wrapping.
func (e *ChartError) Unwrap() error {
	return e.Underlying
}

// Is implements error comparison for errors.Is
func (e *ChartError) Is(target error) bool {
	if chartErr, ok := target.(*ChartError); ok {
		return e.Type == chartErr.Type && e.Operation == chartErr.Operation
	}
	return false
}

// WithContext adds context to an error
func (e *ChartError) WithContext(key, value string) *ChartError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// WithFile adds file context to an error
func (e *ChartError) WithFile(filePath string) *ChartError {
	return e.WithContext("file", filePath)
}

// WithChart adds chart context to an error
func (e *ChartError) WithChart(chartName string) *ChartError {
	return e.WithContext("chart", chartName)
}

// WrapError wraps an existing error with context
func WrapError(err error, errorType ErrorType, operation, message string) *ChartError {
	return &ChartError{
		Type:       errorType,
		Operation:  operation,
		Message:    message,
		Underlying: err,
		Context:    make(map[string]string),
	}
}