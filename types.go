package tlm

import "github.com/rcmaniac25/tlm/logging"

type TLMInitialization struct {
	Logging *logging.TLMLoggingInitialization
}

type TLMBreakdown struct {
	Log logging.TLMLogger
}
