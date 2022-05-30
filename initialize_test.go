package tlm

import "testing"

func TestStartupNil(t *testing.T) {
	_, err := Startup(nil)
	if err == nil {
		t.Fatalf("nil initialiation should should have caused error")
	}
}

func TestStartupEmpty(t *testing.T) {
	inits := new(TLMInitialization)
	_, err := Startup(inits)
	if err != nil {
		//TODO: is this true?
		t.Fatalf("empty initialiation should not cause failure")
	}
}
