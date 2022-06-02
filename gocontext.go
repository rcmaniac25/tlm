package tlm

import "context"

type breakdownKey struct{}

var contextBreakdownKey = breakdownKey{}

func contextBreakdown(ctx context.Context) (TLMBreakdown, bool) {
	value := ctx.Value(contextBreakdownKey)
	breakdown, ok := value.(TLMBreakdown)
	return breakdown, ok
}
func Breakdown(ctx context.Context) TLMBreakdown {
	if breakdown, ok := contextBreakdown(ctx); ok {
		return breakdown
	}
	return TLMBreakdown{}
}

func contextWithStruct(ctx context.Context, breakdown TLMBreakdown) context.Context {
	return context.WithValue(ctx, contextBreakdownKey, breakdown)
}
