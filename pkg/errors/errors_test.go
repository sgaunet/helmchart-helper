package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestChartError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ChartError
		contains []string
	}{
		{
			name: "basic error",
			err: &ChartError{
				Type:      ValidationError,
				Operation: "validate-chart",
				Message:   "chart name is required",
			},
			contains: []string{"operation: validate-chart", "type: validation", "chart name is required"},
		},
		{
			name: "error with underlying cause",
			err: &ChartError{
				Type:       FileSystemError,
				Operation:  "create-file",
				Message:    "failed to create chart file",
				Underlying: errors.New("permission denied"),
			},
			contains: []string{"operation: create-file", "type: filesystem", "caused by: permission denied"},
		},
		{
			name: "error with context",
			err: &ChartError{
				Type:      TemplateError,
				Operation: "process-template",
				Message:   "template parsing failed",
				Context: map[string]string{
					"file":  "deployment.yaml",
					"chart": "my-app",
				},
			},
			contains: []string{"operation: process-template", "type: template", "context:", "file=deployment.yaml", "chart=my-app"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()
			
			for _, expected := range tt.contains {
				if !strings.Contains(errStr, expected) {
					t.Errorf("Error string '%s' does not contain expected substring '%s'", errStr, expected)
				}
			}
		})
	}
}

func TestChartError_Unwrap(t *testing.T) {
	underlying := errors.New("original error")
	chartErr := &ChartError{
		Type:       FileSystemError,
		Operation:  "test",
		Underlying: underlying,
	}

	if chartErr.Unwrap() != underlying {
		t.Errorf("Unwrap() returned %v, expected %v", chartErr.Unwrap(), underlying)
	}
}

func TestChartError_Is(t *testing.T) {
	tests := []struct {
		name   string
		err1   *ChartError
		err2   *ChartError
		expect bool
	}{
		{
			name: "same type and operation",
			err1: &ChartError{Type: ValidationError, Operation: "validate"},
			err2: &ChartError{Type: ValidationError, Operation: "validate"},
			expect: true,
		},
		{
			name: "different type",
			err1: &ChartError{Type: ValidationError, Operation: "validate"},
			err2: &ChartError{Type: FileSystemError, Operation: "validate"},
			expect: false,
		},
		{
			name: "different operation",
			err1: &ChartError{Type: ValidationError, Operation: "validate"},
			err2: &ChartError{Type: ValidationError, Operation: "create"},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err1.Is(tt.err2)
			if result != tt.expect {
				t.Errorf("Is() returned %v, expected %v", result, tt.expect)
			}
		})
	}
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name        string
		constructor func() *ChartError
		expectedType ErrorType
		expectedOp   string
	}{
		{
			name: "NewValidationError",
			constructor: func() *ChartError {
				return NewValidationError("validate-chart", "invalid chart name")
			},
			expectedType: ValidationError,
			expectedOp:   "validate-chart",
		},
		{
			name: "NewFileSystemError",
			constructor: func() *ChartError {
				return NewFileSystemError("create-dir", "directory creation failed", errors.New("permission denied"))
			},
			expectedType: FileSystemError,
			expectedOp:   "create-dir",
		},
		{
			name: "NewTemplateError",
			constructor: func() *ChartError {
				return NewTemplateError("parse-template", "template syntax error", errors.New("unexpected token"))
			},
			expectedType: TemplateError,
			expectedOp:   "parse-template",
		},
		{
			name: "NewConfigurationError",
			constructor: func() *ChartError {
				return NewConfigurationError("load-config", "invalid configuration format")
			},
			expectedType: ConfigurationError,
			expectedOp:   "load-config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor()
			
			if err.Type != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, err.Type)
			}
			
			if err.Operation != tt.expectedOp {
				t.Errorf("Expected operation %s, got %s", tt.expectedOp, err.Operation)
			}
		})
	}
}

func TestChartError_WithContext(t *testing.T) {
	err := NewValidationError("test", "test message")
	
	// Test WithContext
	err = err.WithContext("key1", "value1")
	if err.Context["key1"] != "value1" {
		t.Errorf("Expected context key1=value1, got %s", err.Context["key1"])
	}
	
	// Test WithFile
	err = err.WithFile("/path/to/file.yaml")
	if err.Context["file"] != "/path/to/file.yaml" {
		t.Errorf("Expected file context, got %s", err.Context["file"])
	}
	
	// Test WithChart
	err = err.WithChart("my-chart")
	if err.Context["chart"] != "my-chart" {
		t.Errorf("Expected chart context, got %s", err.Context["chart"])
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	
	wrappedErr := WrapError(originalErr, FileSystemError, "test-operation", "test message")
	
	if wrappedErr.Type != FileSystemError {
		t.Errorf("Expected type %s, got %s", FileSystemError, wrappedErr.Type)
	}
	
	if wrappedErr.Operation != "test-operation" {
		t.Errorf("Expected operation test-operation, got %s", wrappedErr.Operation)
	}
	
	if wrappedErr.Underlying != originalErr {
		t.Errorf("Expected underlying error to be preserved")
	}
}