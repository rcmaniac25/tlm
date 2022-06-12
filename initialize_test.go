package tlm_test

import (
	"testing"

	"github.com/rcmaniac25/tlm"
	"github.com/rcmaniac25/tlm/logging"
)

func TestStartupNil(t *testing.T) {
	_, err := tlm.Startup(nil)
	if err == nil {
		t.Fatal("nil initialization should have caused error")
	}
	if err.Error() != "args cannot be nil" {
		t.Fatalf("Unexpected error message: %s", err.Error())
	}
}

func TestStartupEmpty(t *testing.T) {
	inits := new(tlm.TLMInitialization)
	_, err := tlm.Startup(inits)
	if err == nil {
		t.Fatal("empty initialization should have caused error")
	}
	if err.Error() != "initialization args empty" {
		t.Fatalf("Unexpected error message: %s", err.Error())
	}
}

func TestStartupEmptyLogging(t *testing.T) {
	inits := new(tlm.TLMInitialization)
	inits.Logging = new(logging.TLMLoggingInitialization)
	_, err := tlm.Startup(inits)
	if err == nil {
		t.Fatal("empty logging initialization should have caused error")
	}
	// Error message and logging stuff will be handled by other tests
}
