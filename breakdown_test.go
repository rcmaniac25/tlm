package tlm_test

import (
	"context"
	"testing"

	"github.com/rcmaniac25/tlm"
	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
)

func TestBreakdown(t *testing.T) {
	inits := new(tlm.TLMInitialization)
	inits.Logging = new(logging.TLMLoggingInitialization)

	collector := logging.NewDebugLogCollector()
	collector.SetupInitialization(inits.Logging)

	ctx, _ := tlm.Startup(inits)
	breakdown := tlm.Breakdown(ctx)
	util.AssertEqual(t, breakdown.Ctx, ctx, "context")
	util.AssertNotEqual(t, breakdown.Log, nil, "log")
}

func TestNilBreakdown(t *testing.T) {
	breakdown := tlm.Breakdown(context.Background())
	util.AssertEqual(t, breakdown.Ctx, nil, "context")
	util.AssertEqual(t, breakdown.Log, nil, "log")
}
