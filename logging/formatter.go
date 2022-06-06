package logging

type FormatterType int

const (
	DefaultFormat FormatterType = iota

	TextFormat
	JsonFormat
)

func (g FormatterType) String() string {
	switch g {
	case TextFormat:
		return "text"
	case JsonFormat:
		return "json"
	}
	return ""
}

type Formatter struct {
	Type FormatterType

	// Set the keys used for each log entry. Some defaults available:
	// "-" means the key is skipped (only certain fields supported)
	// "~" means the key is lower-case variable name, without the "Key" suffix
	// ""  means the logger default is used
	TimeKey     string // Can be skipped with "-"
	MessageKey  string
	LevelKey    string
	FunctionKey string

	// Default of time.RFC3339 is used if not set
	TimeFormat string
}
