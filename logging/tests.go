package logging

import (
	"bytes"
	"encoding/json"
)

const (
	LogMessageKey = "lc_log_message"
)

type DebugLogCollector struct {
	buffer bytes.Buffer

	fields map[string]any
}

func NewDebugLogCollector() *DebugLogCollector {
	return new(DebugLogCollector)
}

func (c *DebugLogCollector) SetupInitialization(init *TLMLoggingInitialization) {
	init.Type = LogrusLogType
	init.Output = &c.buffer
	init.Formatter.Type = JsonFormat
	init.Formatter.MessageKey = LogMessageKey
}

func (c *DebugLogCollector) checkDecode() error {
	if len(c.fields) > 0 {
		return nil
	}
	decoder := json.NewDecoder(&c.buffer)
	return decoder.Decode(&c.fields)
}

func (c *DebugLogCollector) GetMessage() string {
	if err := c.checkDecode(); err != nil {
		return err.Error()
	}

	if msg, ok := c.fields[LogMessageKey].(string); ok {
		return msg
	}
	return "Unknown type of log message. Expected string"
}
