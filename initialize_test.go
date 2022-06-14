package tlm_test

import (
	"testing"

	"github.com/rcmaniac25/tlm"
	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
)

func TestStartupNil(t *testing.T) {
	_, err := tlm.Startup(nil)
	util.AssertError(t, err, "nil initialization should have caused error")
	util.AssertEqual(t, err.Error(), "args cannot be nil", "Unexpected error message")
}

func TestStartupEmpty(t *testing.T) {
	inits := new(tlm.TLMInitialization)
	_, err := tlm.Startup(inits)
	util.AssertError(t, err, "empty initialization should have caused error")
	util.AssertEqual(t, err.Error(), "initialization args empty", "Unexpected error message")
}

func TestStartupEmptyLogging(t *testing.T) {
	inits := new(tlm.TLMInitialization)
	inits.Logging = new(logging.TLMLoggingInitialization)
	_, err := tlm.Startup(inits)
	util.AssertError(t, err, "empty logging initialization should have caused error")
	// Error message and logging stuff will be handled by other tests
}

func TestStartupLogging(t *testing.T) {
	inits := new(tlm.TLMInitialization)
	inits.Logging = new(logging.TLMLoggingInitialization)

	collector := logging.NewDebugLogCollector()
	collector.SetupInitialization(inits.Logging)

	ctx, err := tlm.Startup(inits)
	util.AssertNoError(t, err, "logging initialization should not have caused an error")

	tlm.Log(ctx).Info("Hello Tester")
	util.AssertEqual(t, collector.GetNumberLogs(), 1, "log count")
	util.AssertEqual(t, collector.GetMessage(0), "Hello Tester", "log output")
}
