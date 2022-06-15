package logging_test

import (
	"fmt"
	"testing"

	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
)

func TestLogType(t *testing.T) {
	tests := []struct {
		name     string
		value    logging.LogType
		expected string
	}{
		{
			name:     "Unknown",
			value:    -2,
			expected: "unknown",
		},
		{
			name:     "Custom",
			value:    logging.CustomLogType,
			expected: "Custom",
		},
		{
			name:     "Logrus",
			value:    logging.LogrusLogType,
			expected: "Logrus",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stringVersion := fmt.Sprint(tt.value)
			util.AssertEqual(t, stringVersion, tt.expected, "type")
		})
	}
}

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		value    logging.LogLevel
		expected string
	}{
		{
			name:     "Unknown",
			value:    -2,
			expected: "",
		},
		{
			name:     "Default",
			value:    logging.DefaultLevel,
			expected: "",
		},
		{
			name:     "Debug",
			value:    logging.DebugLevel,
			expected: "debug",
		},
		{
			name:     "Info",
			value:    logging.InfoLevel,
			expected: "info",
		},
		{
			name:     "Warn",
			value:    logging.WarnLevel,
			expected: "warn",
		},
		{
			name:     "Error",
			value:    logging.ErrorLevel,
			expected: "error",
		},
		{
			name:     "Panic",
			value:    logging.PanicLevel,
			expected: "panic",
		},
		{
			name:     "Fatal",
			value:    logging.FatalLevel,
			expected: "fatal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stringVersion := fmt.Sprint(tt.value)
			util.AssertEqual(t, stringVersion, tt.expected, "level")
		})
	}
}
