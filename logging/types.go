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

type LogLevel int

const (
	// Default log level is logger specific
	Default LogLevel = iota

	Debug
	Info
	Warn
	Error
	Panic
	Fatal
)

func (t LogLevel) String() string {
	switch t {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	case Panic:
		return "panic"
	case Fatal:
		return "fatal"
	}
	return ""
}

type TLMLoggingInitialization struct {
	Type LogType
	// Only used when Type is Custom
	CustomeType string

	Output io.Writer
	Level  LogLevel
}

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Panicf(format string, args ...any)
	Fatalf(format string, args ...any)

	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Panic(args ...any)
	Fatal(args ...any)

	Debugln(args ...any)
	Infoln(args ...any)
	Warnln(args ...any)
	Errorln(args ...any)
	Panicln(args ...any)
	Fatalln(args ...any)
}

type TLMLogger interface {
	Logger

	Context() context.Context
}
