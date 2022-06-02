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

func (n *nullLoggerType) Info(msg string) {}
func (s *selfReferentialLogger) Info(msg string) {
	s.LoggerImpl.Info(msg)
}
