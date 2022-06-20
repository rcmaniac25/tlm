package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"time"
)

const (
	LogMessageKey = "lc_log_message"
	LogLevelKey   = "lc_log_level"
	LogTimeKey    = "lc_log_time"
)

type DebugLogCollector struct {
	buffer bytes.Buffer

	logs      []map[string]any
	exitcodes []int
}

func NewDebugLogCollector() *DebugLogCollector {
	return new(DebugLogCollector)
}

func (c *DebugLogCollector) SetupInitialization(init *TLMLoggingInitialization) {
	init.Type = LogrusLogType
	init.Output = &c.buffer
	init.Formatter.Type = JsonFormat
	init.Formatter.MessageKey = LogMessageKey
	init.Formatter.LevelKey = LogLevelKey
	init.Formatter.TimeKey = LogTimeKey
	init.Formatter.TimeFormat = time.RFC3339Nano
}

func (c *DebugLogCollector) OnExitCode(exitcode int) {
	c.exitcodes = append(c.exitcodes, exitcode)
}

func (c *DebugLogCollector) populateLogs() error {
	decoder := json.NewDecoder(&c.buffer)
	for {
		var fields map[string]any
		err := decoder.Decode(&fields)
		if err == nil {
			c.logs = append(c.logs, fields)
		}
		if err != nil && err == io.EOF {
			// Reached the end
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (c *DebugLogCollector) Clear() {
	c.populateLogs()
	c.logs = make([]map[string]any, 0)
	c.exitcodes = make([]int, 0)
}

func (c *DebugLogCollector) GetNumberLogs() int {
	if err := c.populateLogs(); err != nil {
		return -1
	}
	return len(c.logs)
}

func (c *DebugLogCollector) getLogField(logIndex int, field string) (any, error) {
	if logIndex < 0 {
		return nil, errors.New("<dev> log index must be >= 0")
	}

	if logIndex >= len(c.logs) {
		if err := c.populateLogs(); err != nil {
			return nil, err
		}
	}

	if logIndex >= len(c.logs) {
		return nil, errors.New("<dev> log index exceeds number of logs received")
	}

	return c.logs[logIndex][field], nil
}

func (c *DebugLogCollector) GetMessage(logIndex int) string {
	field, err := c.getLogField(logIndex, LogMessageKey)
	if err != nil {
		return err.Error()
	}
	if msg, ok := field.(string); ok {
		return msg
	}
	return "Unknown type of log message. Expected string"
}

func (c *DebugLogCollector) GetLogLevel(logIndex int) LogLevel {
	field, err := c.getLogField(logIndex, LogLevelKey)
	if err != nil {
		return -1
	}
	level, ok := field.(string)
	if !ok {
		return -2
	}
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	// The next 2 levels could potentially get lost due to the log level itself
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		// Not actually default, but gives a good way to go "no idea"
		return DefaultLevel
	}
}

func (c *DebugLogCollector) GetTime(logIndex int) time.Time {
	field, err := c.getLogField(logIndex, LogTimeKey)
	if err != nil {
		return time.Time{}
	}
	if timeStr, ok := field.(string); ok {
		if t, err := time.Parse(time.RFC3339Nano, timeStr); err == nil {
			return t
		}
	}
	return time.Time{}
}

func (c *DebugLogCollector) GetField(logIndex int, field string) (any, bool) {
	switch field {
	case LogMessageKey, LogLevelKey, LogTimeKey:
		return nil, false
	}
	if f, err := c.getLogField(logIndex, field); err == nil {
		return f, true
	}
	return nil, false
}

func (c *DebugLogCollector) GetFieldFunc(logIndex int, field string) func() (any, bool) {
	return func() (any, bool) {
		return c.GetField(logIndex, field)
	}
}

func (c *DebugLogCollector) GetFatalExitcode(logIndex int) (int, bool) {
	if logIndex < 0 {
		return -1, false
	}

	if logIndex >= len(c.logs) {
		if err := c.populateLogs(); err != nil {
			return -2, false
		}
	}

	if logIndex >= len(c.logs) {
		return -3, false
	}

	// This attempts to sync fatal function calls to error codes
	type exitcodesync struct {
		valid bool
		value int
	}
	errcodeIndex := 0
	var codes []exitcodesync
	for i, log := range c.logs {
		if i > logIndex {
			// No need to keep going if we already reached the index
			break
		}
		if errcodeIndex >= len(c.exitcodes) {
			// No need to keep going if there are no more codes to populate
			break
		}
		if log[LogLevelKey] != "fatal" {
			codes = append(codes, exitcodesync{
				valid: false,
				value: 0,
			})
			continue
		}

		codes = append(codes, exitcodesync{
			valid: true,
			value: c.exitcodes[errcodeIndex],
		})
		errcodeIndex++
	}
	if logIndex >= len(codes) {
		return -4, false
	}
	return codes[logIndex].value, codes[logIndex].valid
}

func (c *DebugLogCollector) GetFatalExitcodeFunc(logIndex int) func() (any, bool) {
	return func() (any, bool) {
		return c.GetFatalExitcode(logIndex)
	}
}
