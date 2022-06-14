package tlm

import "context"

type breakdownKey struct{}

var contextBreakdownKey = breakdownKey{}

func contextBreakdown(ctx context.Context) (TLMBreakdown, bool) {
	value := ctx.Value(contextBreakdownKey)
	breakdown, ok := value.(TLMBreakdown)
	if ok {
		breakdown.Ctx = ctx
	}
	return breakdown, ok
}

func contextWithStruct(ctx context.Context, breakdown TLMBreakdown) context.Context {
	return context.WithValue(ctx, contextBreakdownKey, breakdown)
}
