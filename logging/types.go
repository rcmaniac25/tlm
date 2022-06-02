package logging

import (
	"context"
	"io"
)

type CustomeLoggerInitializationFunc func(args *TLMLoggingInitialization) (Logger, error)

type LogType int

const (
	// Requires setting CustomeType during initialization
	Custom LogType = iota

	Logrus
)

func (t LogType) String() string {
	switch t {
	case Custom:
		return "Custom"
	case Logrus:
		return "Logrus"
	}
	return "unknown"
}

type TLMLoggingInitialization struct {
	Type LogType
	// Only used when Type is Custom
	CustomeType string

	Output io.Writer
}

type Logger interface {
	Info(msg string)
}

type TLMLogger interface {
	Logger

	Context() context.Context
}
