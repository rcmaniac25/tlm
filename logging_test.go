package tlm_test

import (
	"context"
	"testing"

	"github.com/rcmaniac25/tlm"
	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
)

func TestLogging(t *testing.T) {
	inits := new(tlm.TLMInitialization)
	inits.Logging = new(logging.TLMLoggingInitialization)

	collector := logging.NewDebugLogCollector()
	collector.SetupInitialization(inits.Logging)

	ctx, _ := tlm.Startup(inits)
	log := tlm.Log(ctx)
	util.AssertNotEqual(t, log, nil, "log")

	log.Info("Hello Tester")
	util.AssertEqual(t, collector.GetMessage(), "Hello Tester", "log output")
}

func TestLoggingNoInit(t *testing.T) {
	ctx := context.Background()
	log := tlm.Log(ctx)
	util.AssertNotEqual(t, log, nil, "log")
	util.AssertNoPanic(t, func() {
		log.Info("Hello Tester")
	}, "nil log")
}
