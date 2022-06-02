package tlm_test

import (
	"testing"

	"github.com/rcmaniac25/tlm"
)

func TestStartupNil(t *testing.T) {
	_, err := tlm.Startup(nil)
	if err == nil {
		t.Fatalf("nil initialiation should should have caused error")
	}
}

func TestStartupEmpty(t *testing.T) {
	inits := new(tlm.TLMInitialization)
	ctx, err := tlm.Startup(inits)
	if err != nil {
		//TODO: is this true?
		t.Fatalf("empty initialiation should not cause failure")
	}
	tlm.Log(ctx).Info("Test") //TMP
}
