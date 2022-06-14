package tlm

import (
	"context"

	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
)

type TLMBreakdown struct {
	Log logging.TLMLogger
	Ctx context.Context
}

func Breakdown(ctx context.Context) TLMBreakdown {
	if breakdown, ok := contextBreakdown(ctx); ok {
		return breakdown
	}
	return TLMBreakdown{}
}

// This exists to allow updating a breakdown without exposing it's internals and causing a import cycle
type tlmBreakdownContextWrapper struct {
	Ctx context.Context
}

func (w tlmBreakdownContextWrapper) GetContext() context.Context {
	return w.Ctx
}

func (w tlmBreakdownContextWrapper) UpdateLogger(logger logging.TLMLogger) util.ContextWrapper {
	if breakdown, ok := contextBreakdown(w.Ctx); ok {
		breakdown.Log = logger
		return &tlmBreakdownContextWrapper{Ctx: contextWithStruct(w.Ctx, breakdown)}
	}
	return w
}
