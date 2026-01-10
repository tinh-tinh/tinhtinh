package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLevelName(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelDebug, "debug"},
		{LevelInfo, "info"},
		{LevelWarn, "warn"},
		{LevelError, "error"},
		{LevelFatal, "fatal"},
		{Level(99), ""},
	}

	for _, tt := range tests {
		result := GetLevelName(tt.level)
		require.Equal(t, tt.expected, result)
	}
}

func TestExtractAllContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single variable",
			input:    "${method}",
			expected: []string{"method"},
		},
		{
			name:     "multiple variables",
			input:    "${ip} - ${method} ${path} ${status}",
			expected: []string{"ip", "method", "path", "status"},
		},
		{
			name:     "no variables",
			input:    "plain text without variables",
			expected: nil,
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "complex format",
			input:    "${ip}:${date} - '${method} ${path} ${http-version}' ${status} ${content-length}",
			expected: []string{"ip", "date", "method", "path", "http-version", "status", "content-length"},
		},
		{
			name:     "variable with hyphen",
			input:    "${user-agent}",
			expected: []string{"user-agent"},
		},
		{
			name:     "nested braces",
			input:    "${{nested}}",
			expected: []string{"{nested"},
		},
		{
			name:     "incomplete variable",
			input:    "${incomplete",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAllContent(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
