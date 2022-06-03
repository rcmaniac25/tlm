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
}
