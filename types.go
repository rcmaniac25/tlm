package tlm

// LogType
type LogType int

const (
	Logrus LogType = iota
)

func (t LogType) String() string {
	//TODO: expand
	if t != Logrus {
		return "unknown"
	}

	return "Logrus"
}

// TLMLoggingInitialization
type TLMLoggingInitialization struct {
	Type LogType
}

// TLMInitialization
type TLMInitialization struct {
	Logging *TLMLoggingInitialization
}
