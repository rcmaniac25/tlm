package tlm

import (
	"context"

	"github.com/rcmaniac25/tlm/logging"
)

func Log(ctx context.Context) logging.Logger {
	if breakdown, ok := contextBreakdown(ctx); ok && breakdown.Log != nil {
		return breakdown.Log
	}
	return &logging.NullLogger
}
