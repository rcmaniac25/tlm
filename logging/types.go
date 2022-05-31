package logging

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

type TLMLoggingInitialization struct {
	Type LogType
}

type Logger struct {
	Impl interface{} //TODO
}
