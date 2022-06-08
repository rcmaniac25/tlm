package logging

import (
	"context"
	"io"
)

type CustomeLoggerInitializationFunc func(args *TLMLoggingInitialization) (Logger, error)

type LogType int

const (
	// Requires setting CustomeType during initialization
	CustomLogType LogType = iota

	LogrusLogType
)

func (t LogType) String() string {
	switch t {
	case CustomLogType:
		return "Custom"
	case LogrusLogType:
		return "Logrus"
	}
	return "unknown"
}

type LogLevel int

const (
	// Default log level is logger specific
	DefaultLevel LogLevel = iota

	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

func (t LogLevel) String() string {
	switch t {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	}
	return ""
}

type TLMLoggingInitialization struct {
	Type LogType
	// Only used when Type is Custom
	CustomeType string

	Output    io.Writer
	Level     LogLevel
	Formatter Formatter

	//TODO: logger specific variables
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

	WithField(key string, value any) Logger
	//TODO: WithFields

	//TODO: PanicOnlyDebugMode //XXX don't actually implement this. This should be a value passed into WithField(s) and if
}

type TLMLogger interface {
	Logger

	Context() context.Context
}
