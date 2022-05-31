package tlm

import (
	"context"
	"errors"

	log "github.com/rcmaniac25/tlm/logging"
)

func Startup(init *TLMInitialization) (context.Context, error) {
	if init == nil {
		return context.Background(), errors.New("initialization values cannot be nil")
	}

	var breakdown TLMBreakdown

	//TODO: tracing

	logger, err := log.InitLogging(init.Logging)
	if err != nil {
		return context.Background(), err
	}
	breakdown.Log = logger

	//TODO: metrics

	return contextWithStruct(context.Background(), breakdown), nil
}
