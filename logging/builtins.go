package logging

import (
	"context"
	"errors"
)

type nullLoggerType struct{}

var NullLogger = nullLoggerType{}

// This basically creates a reference loop. We want to be able to chain log, metric, etc. updates.
// In order to do that, we need access to TLM's breakdown... which is stored in the context.
// So we need to store the context that contains the logger, in the logger
type selfReferentialLogger struct {
	TLMContext context.Context
	LoggerImpl Logger
}

// Not a selfReferentialLogger func on purpose, because it won't be known until exported
func SetTLMLoggerContext(logger TLMLogger, ctx context.Context) error {
	if selfRefLogger, ok := logger.(*selfReferentialLogger); ok {
		selfRefLogger.TLMContext = ctx
		return nil
	}
	return errors.New("unsupported logger")
}

func (s *selfReferentialLogger) Context() context.Context {
	return s.TLMContext
}

// All the builtin functions

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

func (n *nullLoggerType) Panicf(format string, args ...any) {}
func (n *nullLoggerType) Panic(args ...any)                 {}
func (n *nullLoggerType) Panicln(args ...any)               {}
func (s *selfReferentialLogger) Panicf(format string, args ...any) {
	s.LoggerImpl.Panicf(format, args...)
}
func (s *selfReferentialLogger) Panic(args ...any) {
	s.LoggerImpl.Panic(args...)
}
func (s *selfReferentialLogger) Panicln(args ...any) {
	s.LoggerImpl.Panicln(args...)
}

func (n *nullLoggerType) Fatalf(format string, args ...any) {}
func (n *nullLoggerType) Fatal(args ...any)                 {}
func (n *nullLoggerType) Fatalln(args ...any)               {}
func (s *selfReferentialLogger) Fatalf(format string, args ...any) {
	s.LoggerImpl.Fatalf(format, args...)
}
func (s *selfReferentialLogger) Fatal(args ...any) {
	s.LoggerImpl.Fatal(args...)
}
func (s *selfReferentialLogger) Fatalln(args ...any) {
	s.LoggerImpl.Fatalln(args...)
}
