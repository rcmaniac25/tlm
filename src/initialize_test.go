package tlm

import "testing"

func TestStartupEmpty(t *testing.T) {
	inits := new(TLMInitialization)
	_, err := Startup(inits)
	if err != nil {
		//TODO: is this true?
		t.Fatalf("empty initialiation should not cause failue")
	}
}
