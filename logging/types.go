package logging

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

type Logger interface {
	Info(msg string)
}

//TODO: probably move this to it's own file
type nullLoggerType struct{}

var NullLogger = nullLoggerType{}

func (n *nullLoggerType) Info(msg string) {}
