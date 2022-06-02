package logging

type CustomeLoggerInitializationFunc func(args *TLMLoggingInitialization) (Logger, error)

type LogType int

const (
	// Requires setting CustomeType during initialization
	Custom LogType = iota

	Logrus
)

func (t LogType) String() string {
	switch t {
	case Custom:
		return "Custom"
	case Logrus:
		return "Logrus"
	}
	return "unknown"
}

type TLMLoggingInitialization struct {
	Type LogType
	// Only used when Type is Custom
	CustomeType string
}

type Logger interface {
	Info(msg string)
}

//TODO: probably move this to it's own file
type nullLoggerType struct{}

var NullLogger = nullLoggerType{}

func (n *nullLoggerType) Info(msg string) {}
