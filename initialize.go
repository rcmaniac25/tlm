package tlm

import (
	"context"
	"errors"

	"github.com/rcmaniac25/tlm/logging"
)

func Startup(args *TLMInitialization) (context.Context, error) {
	return StartupContext(context.Background(), args)
}

func StartupContext(ctx context.Context, args *TLMInitialization) (context.Context, error) {
	if args == nil {
		return context.Background(), errors.New("args cannot be nil")
	}

	var breakdown TLMBreakdown

	//TODO: tracing

	logger, err := logging.InitLogging(args.Logging)
	if err != nil {
		return context.Background(), err
	}
	breakdown.Log = logger

	//TODO: metrics

	tlmCtx := contextWithStruct(ctx, breakdown)

	type setContext interface {
		SetContext(ctx context.Context)
	}
	if breakdown.Log != nil {
		if setCtx, ok := breakdown.Log.(setContext); ok {
			setCtx.SetContext(tlmCtx)
		} else {
			return nil, errors.New("internal error: unknown logger")
		}
	}

	return tlmCtx, nil
}
