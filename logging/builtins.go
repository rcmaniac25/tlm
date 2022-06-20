package logging

import (
	"context"
	"os"

	"github.com/rcmaniac25/tlm/util"
)

type nullLoggerType struct{}

var NullLogger = nullLoggerType{}

// This basically creates a reference loop. We want to be able to chain log, metric, etc. updates.
// In order to do that, we need access to TLM's breakdown... which is stored in the context.
// So we need to store the context that contains the logger, in the logger
type selfReferentialLogger struct {
	TLMContext util.ContextWrapper
	LoggerImpl Logger
}

func (s *selfReferentialLogger) SetContextWrapper(ctx util.ContextWrapper) {
	s.TLMContext = ctx
}

func (s *selfReferentialLogger) Context() context.Context {
	return s.TLMContext.GetContext()
}

func (s *selfReferentialLogger) updateLogger(update func() Logger) Logger {
	type UpdateLogger interface {
		UpdateLogger(logger TLMLogger) util.ContextWrapper
	}

	refLogger := &selfReferentialLogger{
		LoggerImpl: update(),
	}
	if updateLogger, ok := s.TLMContext.(UpdateLogger); ok {
		refLogger.TLMContext = updateLogger.UpdateLogger(refLogger)
		return refLogger
	}
	return s // Simply ignore the field since we got an invalid type...
}

func (s *selfReferentialLogger) TestingSetFatalExitFunction(exitHandler func(int)) bool {
	type InternalTestingExitHandler interface {
		TestingSetFatalExitFunction(exitHandler func(int)) bool
	}
	if v, ok := s.LoggerImpl.(InternalTestingExitHandler); ok {
		return v.TestingSetFatalExitFunction(exitHandler)
	}
	return false
}

// All the builtin functions

func (n *nullLoggerType) WithField(key string, value any) Logger {
	return n
}
func (s *selfReferentialLogger) WithField(key string, value any) Logger {
	return s.updateLogger(func() Logger {
		return s.LoggerImpl.WithField(key, value)
	})
}

func (n *nullLoggerType) WithFields(fields util.Fields) Logger {
	return n
}
func (s *selfReferentialLogger) WithFields(fields util.Fields) Logger {
	return s.updateLogger(func() Logger {
		return s.LoggerImpl.WithFields(fields)
	})
}

func (n *nullLoggerType) Debugf(format string, args ...any) {}
func (n *nullLoggerType) Debug(args ...any)                 {}
func (n *nullLoggerType) Debugln(args ...any)               {}
func (s *selfReferentialLogger) Debugf(format string, args ...any) {
	s.LoggerImpl.Debugf(format, args...)
}
func (s *selfReferentialLogger) Debug(args ...any) {
	s.LoggerImpl.Debug(args...)
}
func (s *selfReferentialLogger) Debugln(args ...any) {
	s.LoggerImpl.Debugln(args...)
}

func (n *nullLoggerType) Infof(format string, args ...any) {}
func (n *nullLoggerType) Info(args ...any)                 {}
func (n *nullLoggerType) Infoln(args ...any)               {}
func (s *selfReferentialLogger) Infof(format string, args ...any) {
	s.LoggerImpl.Infof(format, args...)
}
func (s *selfReferentialLogger) Info(args ...any) {
	s.LoggerImpl.Info(args...)
}
func (s *selfReferentialLogger) Infoln(args ...any) {
	s.LoggerImpl.Infoln(args...)
}

func (n *nullLoggerType) Warnf(format string, args ...any) {}
func (n *nullLoggerType) Warn(args ...any)                 {}
func (n *nullLoggerType) Warnln(args ...any)               {}
func (s *selfReferentialLogger) Warnf(format string, args ...any) {
	s.LoggerImpl.Warnf(format, args...)
}
func (s *selfReferentialLogger) Warn(args ...any) {
	s.LoggerImpl.Warn(args...)
}
func (s *selfReferentialLogger) Warnln(args ...any) {
	s.LoggerImpl.Warnln(args...)
}

func (n *nullLoggerType) Errorf(format string, args ...any) {}
func (n *nullLoggerType) Error(args ...any)                 {}
func (n *nullLoggerType) Errorln(args ...any)               {}
func (s *selfReferentialLogger) Errorf(format string, args ...any) {
	s.LoggerImpl.Errorf(format, args...)
}
func (s *selfReferentialLogger) Error(args ...any) {
	s.LoggerImpl.Error(args...)
}
func (s *selfReferentialLogger) Errorln(args ...any) {
	s.LoggerImpl.Errorln(args...)
}

func (n *nullLoggerType) Panicf(format string, args ...any) { panic("Panicf") }
func (n *nullLoggerType) Panic(args ...any)                 { panic("Panic") }
func (n *nullLoggerType) Panicln(args ...any)               { panic("Panicln") }
func (s *selfReferentialLogger) Panicf(format string, args ...any) {
	s.LoggerImpl.Panicf(format, args...)
}
func (s *selfReferentialLogger) Panic(args ...any) {
	s.LoggerImpl.Panic(args...)
}
func (s *selfReferentialLogger) Panicln(args ...any) {
	s.LoggerImpl.Panicln(args...)
}

func (n *nullLoggerType) Fatalf(format string, args ...any) { os.Exit(1) }
func (n *nullLoggerType) Fatal(args ...any)                 { os.Exit(1) }
func (n *nullLoggerType) Fatalln(args ...any)               { os.Exit(1) }
func (s *selfReferentialLogger) Fatalf(format string, args ...any) {
	s.LoggerImpl.Fatalf(format, args...)
}
func (s *selfReferentialLogger) Fatal(args ...any) {
	s.LoggerImpl.Fatal(args...)
}
func (s *selfReferentialLogger) Fatalln(args ...any) {
	s.LoggerImpl.Fatalln(args...)
}
