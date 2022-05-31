package tlm

import "context"

type breakdownKey struct{}

var contextBreakdownKey = breakdownKey{}

func Breakdown(ctx context.Context) TLMBreakdown {
	value := ctx.Value(contextBreakdownKey)
	if breakdown, ok := value.(TLMBreakdown); ok {
		return breakdown
	}
	return TLMBreakdown{}
}

func contextWithStruct(ctx context.Context, breakdown TLMBreakdown) context.Context {
	return context.WithValue(ctx, contextBreakdownKey, breakdown)
}