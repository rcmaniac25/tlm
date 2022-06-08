package tlm

import (
	"context"
	"errors"

	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
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

	/* Result is effectivly:
	ContextWrapper {
		context.Context {
			TLMBreakdown {
				...
				context.Context // set when getting the breakdown value out of the context
			}
		}
	}
	*/
	tlmCtxWrapper := tlmBreakdownContextWrapper{
		Ctx: contextWithStruct(ctx, breakdown),
	}

	type setContext interface {
		SetContextWrapper(ctx util.ContextWrapper)
	}
	if breakdown.Log != nil {
		if setCtx, ok := breakdown.Log.(setContext); ok {
			setCtx.SetContextWrapper(tlmCtxWrapper)
		} else {
			return nil, errors.New("internal error: unknown logger")
		}
	}

	return tlmCtxWrapper.GetContext(), nil
}
