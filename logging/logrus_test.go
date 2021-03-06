package logging_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
	"github.com/sirupsen/logrus"
)

// Setting output and detailed logging tests are handled in builtins_test

func createLoggerExt(args *logging.TLMLoggingInitialization, exitHandler func(int)) (logging.Logger, *bytes.Buffer) {
	output := new(bytes.Buffer)
	args.Output = output
	logger, err := logging.InitLogrus(args)
	if err != nil {
		panic(err.Error())
	}
	if exitHandler != nil {
		//TODO: figure out why this doesn't work...
		if _, ok := logger.(*logging.LogrusImpl); !ok {
			panic("Got the wrong logger type?")
		}
		type InternalTestingExitHandler interface {
			testExitFunc(exitHandler func(int)) bool
		}
		if v, ok := logger.(InternalTestingExitHandler); ok {
			if !v.testExitFunc(exitHandler) {
				panic("Logrus exit handler not set")
			}
		} else {
			panic("Could not get test exit handler")
		}
	}
	return logger, output
}

func createLogger(args *logging.TLMLoggingInitialization) (logging.Logger, *bytes.Buffer) {
	return createLoggerExt(args, nil)
}

func TestBasic(t *testing.T) {
	logArgs := new(logging.TLMLoggingInitialization)
	logger, buffer := createLogger(logArgs)

	util.AssertNotEqual(t, logger, nil, "logger exists")

	util.AssertEqual(t, buffer.Len(), 0, "buffer length")
	logger.Info("Hello")
	util.AssertNotEqual(t, buffer.Len(), 0, "buffer length")

	util.AssertContains(t, buffer.String(), "Hello", "contents")
}

func TestLevels(t *testing.T) {
	type args struct {
		level       logging.LogLevel
		ignoreCase  func(logging.Logger)
		levelCase   func(logging.Logger)
		expectPanic bool
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "Debug",
			args: args{
				level:      logging.DebugLevel,
				ignoreCase: func(_ logging.Logger) {}, // Nothing that debug wouldn't cover...
				levelCase:  func(logger logging.Logger) { logger.Debug("DebugSuccess") },
			},
			expected: "DebugSuccess",
		},
		{
			name: "Info",
			args: args{
				level:      logging.InfoLevel,
				ignoreCase: func(logger logging.Logger) { logger.Debug("InfoFail") },
				levelCase:  func(logger logging.Logger) { logger.Info("InfoSuccess") },
			},
			expected: "InfoSuccess",
		},
		{
			name: "Warn",
			args: args{
				level:      logging.WarnLevel,
				ignoreCase: func(logger logging.Logger) { logger.Info("WarnFail") },
				levelCase:  func(logger logging.Logger) { logger.Warn("WarnSuccess") },
			},
			expected: "WarnSuccess",
		},
		{
			name: "Error",
			args: args{
				level:      logging.ErrorLevel,
				ignoreCase: func(logger logging.Logger) { logger.Warn("ErrorFail") },
				levelCase:  func(logger logging.Logger) { logger.Error("ErrorSuccess") },
			},
			expected: "ErrorSuccess",
		},
		{
			name: "Panic",
			args: args{
				level:       logging.PanicLevel,
				expectPanic: true,
				ignoreCase:  func(logger logging.Logger) { logger.Error("PanicFail") },
				levelCase:   func(logger logging.Logger) { logger.Panic("PanicSuccess") },
			},
			expected: "PanicSuccess",
		},
		/* TODO: Golang language spec (and my search foo) are not giving info on how to cast certain types. So I end up with this weird mishmash of "how do I get to a function that I don't wnat to expose public?"
		 * Note: this is why "internal" or "package" modifiers exist in some languages. Makes it possible to test code with a lot less lines of code...
		{
			name: "Fatal",
			args: args{
				level:       logging.FatalLevel,
				expectPanic: true,
				ignoreCase:  func(logger logging.Logger) { logger.Error("FatalFail") },
				levelCase:   func(logger logging.Logger) { logger.Fatal("FatalSuccess") },
			},
			expected: "FatalSuccess",
		},
		*/
	}
	for _, level := range tests {
		t.Run(level.name, func(t *testing.T) {
			logArgs := new(logging.TLMLoggingInitialization)
			logArgs.Level = level.args.level
			//TODO: once exit handler can be set, uncomment this: logger, buffer := createLoggerExt(logArgs, func(_ int) {}) // We don't want fatal log calls to cause problems
			logger, buffer := createLogger(logArgs)

			util.AssertEqual(t, buffer.Len(), 0, "buffer length")
			level.args.ignoreCase(logger)
			util.AssertEqual(t, buffer.Len(), 0, "buffer length")

			if level.args.expectPanic {
				util.AssertPanic(t, func() {
					level.args.levelCase(logger)
				}, "panic")
			} else {
				util.AssertNoPanic(t, func() {
					level.args.levelCase(logger)
				}, "no panic")
			}
			util.AssertNotEqual(t, buffer.Len(), 0, "buffer length")

			util.AssertContains(t, buffer.String(), level.expected, "contents")
		})
	}
}

const LogTimePlaceholder = "logtimeplaceholder"

func AssertTime(t *testing.T, timeKey, timeFormat, timeStr, log string, logTime time.Time) {
	if timeKey == "-" {
		util.AssertNotContains(t, log, fmt.Sprintf("%s=", logrus.FieldKeyTime), "no time marker")
		return
	}
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}
	parsedTime, err := time.Parse(timeFormat, timeStr)
	util.AssertNoError(t, err, "parse time")

	// Due to format inaccuracies, reparse logTime
	logTimeStr := logTime.Format(timeFormat)
	parsedLogTime, err := time.Parse(timeFormat, logTimeStr)
	util.AssertNoError(t, err, "parse log time")

	dur := parsedTime.Sub(parsedLogTime)
	util.AssertEqual(t, dur >= time.Duration(0), true, "log time is logical")
	util.AssertEqual(t, dur < (time.Duration(50)*time.Millisecond), true, "log is within threshold")
}

func SplitQuotedString(log string) []string {
	parts := make([]string, 0)

	inQuotedSting := false

	var b strings.Builder
	for _, c := range log {
		if !inQuotedSting {
			if c == ' ' {
				parts = append(parts, b.String())
				b.Reset()
				continue
			}
			b.WriteRune(c)
			if c == '"' {
				inQuotedSting = true
			}
			continue
		}

		b.WriteRune(c)
		if c == '"' {
			inQuotedSting = false
		}
	}
	if b.Len() > 0 {
		parts = append(parts, b.String())
	}
	if inQuotedSting {
		// Sanity
		return make([]string, 0)
	}
	return parts
}

func SplitAndOrderOutput(log string, format logging.Formatter) (newLog string, logTime string) {
	logMap := make(map[string]any)

	logTime = ""
	timeKey := format.TimeKey
	if timeKey == "" {
		timeKey = logrus.FieldKeyTime
	}
	if timeKey == "~" {
		timeKey = "time"
	}

	switch format.Type {
	case logging.TextFormat:
		if log[len(log)-1] != '\n' {
			// Expected it to end with a newline, which the text formatter does
			return "", ""
		}
		parts := SplitQuotedString(log[:len(log)-1])

		if timeKey != "-" {
			idx := -1
			for i := 0; i < len(parts); i++ {
				if strings.Split(parts[i], "=")[0] == timeKey {
					idx = i
					break
				}
			}
			if idx < 0 {
				// Expected to find time, and didn't
				return "", ""
			}
			logTime = strings.Split(parts[idx], "=")[1]
			logTime = logTime[1 : len(logTime)-1] // remove the quotes
			parts[idx] = fmt.Sprintf("%s=\"%s\"", timeKey, LogTimePlaceholder)
		}

		sort.Slice(parts, func(i, j int) bool {
			ikey := parts[i][:strings.IndexRune(parts[i], '=')]
			jkey := parts[j][:strings.IndexRune(parts[j], '=')]
			return ikey < jkey
		})

		return fmt.Sprintf("%s\n", strings.Join(parts, " ")), logTime
	case logging.JsonFormat:
		if err := json.Unmarshal([]byte(log), &logMap); err != nil {
			return err.Error(), ""
		}
		if timeKey != "-" {
			value, ok := logMap[timeKey]
			if !ok {
				return "can't find time", ""
			}
			strValue, ok := value.(string)
			if !ok {
				return "can't get time string", ""
			}

			logTime = strValue
			logMap[timeKey] = LogTimePlaceholder
		}
		result, err := json.Marshal(logMap)
		if err != nil {
			return err.Error(), ""
		}
		return fmt.Sprintf("%s\n", string(result)), logTime
	}
	return "", ""
}

func GetFileAndFunctionPlaceholder(functionName string) func(string) string {
	// This replaces file and function names with placeholders because a change to a log file or some other thing could break the test
	replaceText := func(log, fieldName, placeholder string) string {
		field := fmt.Sprintf("%s=", fieldName)
		idx := strings.Index(log, field)
		if idx >= 0 {
			hasQuotes := log[idx+len(field)] == '"'
			old := ""
			if hasQuotes {
				s := idx + len(field) + 1
				e := strings.Index(log[s:], "\"") + s
				old = log[s:e]
			} else {
				s := idx + len(field)
				e := strings.Index(log[s:], " ") + s
				old = log[s:e]
			}
			return strings.Replace(log, old, placeholder, 1)
		}
		return log
	}
	replaceJson := func(log, fieldName, placeholder string) string {
		field := fmt.Sprintf("\"%s\":\"", fieldName)
		idx := strings.Index(log, field)
		if idx >= 0 {
			s := idx + len(field)
			e := strings.Index(log[s:], "\"") + s
			old := log[s:e]
			return strings.Replace(log, old, placeholder, 1)
		}
		return log
	}

	return func(log string) string {
		// File
		log = replaceText(log, "file", "fileplaceholder")
		log = replaceJson(log, "file", "fileplaceholder")

		// Function
		log = replaceText(log, functionName, "funcplaceholder")
		log = replaceJson(log, functionName, "funcplaceholder")

		return log
	}
}

func TestFormats(t *testing.T) {
	tests := []struct {
		name        string
		args        logging.Formatter
		postLogFunc func(string) string
		expected    map[logging.FormatterType]string
	}{
		{
			name: "Default",
			args: logging.Formatter{},
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "level=info msg=\"Hello World\" testInt=128 testStr=hello time=\"logtimeplaceholder\"\n",
				logging.JsonFormat: "{\"level\":\"info\",\"msg\":\"Hello World\",\"testInt\":128,\"testStr\":\"hello\",\"time\":\"logtimeplaceholder\"}\n",
			},
		},
		{
			name: "Special Time",
			args: logging.Formatter{
				TimeKey: "bobsTime",
			},
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "bobsTime=\"logtimeplaceholder\" level=info msg=\"Hello World\" testInt=128 testStr=hello\n",
				logging.JsonFormat: "{\"bobsTime\":\"logtimeplaceholder\",\"level\":\"info\",\"msg\":\"Hello World\",\"testInt\":128,\"testStr\":\"hello\"}\n",
			},
		},
		{
			name: "No Time",
			args: logging.Formatter{
				TimeKey: "-",
			},
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "level=info msg=\"Hello World\" testInt=128 testStr=hello\n",
				logging.JsonFormat: "{\"level\":\"info\",\"msg\":\"Hello World\",\"testInt\":128,\"testStr\":\"hello\"}\n",
			},
		},
		{
			name: "Time Key",
			args: logging.Formatter{
				TimeKey: "~",
			},
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "level=info msg=\"Hello World\" testInt=128 testStr=hello time=\"logtimeplaceholder\"\n",
				logging.JsonFormat: "{\"level\":\"info\",\"msg\":\"Hello World\",\"testInt\":128,\"testStr\":\"hello\",\"time\":\"logtimeplaceholder\"}\n",
			},
		},
		{
			name: "Time Format",
			args: logging.Formatter{
				TimeFormat: time.RFC822, // The testing compares parsed dates, so there isn't really a reason to ensure the format was correct, as it would just error when parsing
			},
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "level=info msg=\"Hello World\" testInt=128 testStr=hello time=\"logtimeplaceholder\"\n",
				logging.JsonFormat: "{\"level\":\"info\",\"msg\":\"Hello World\",\"testInt\":128,\"testStr\":\"hello\",\"time\":\"logtimeplaceholder\"}\n",
			},
		},
		{
			name: "Special Message",
			args: logging.Formatter{
				MessageKey: "hearsey",
			},
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "hearsey=\"Hello World\" level=info testInt=128 testStr=hello time=\"logtimeplaceholder\"\n",
				logging.JsonFormat: "{\"hearsey\":\"Hello World\",\"level\":\"info\",\"testInt\":128,\"testStr\":\"hello\",\"time\":\"logtimeplaceholder\"}\n",
			},
		},
		{
			name: "Message Key",
			args: logging.Formatter{
				MessageKey: "~",
			},
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "level=info message=\"Hello World\" testInt=128 testStr=hello time=\"logtimeplaceholder\"\n",
				logging.JsonFormat: "{\"level\":\"info\",\"message\":\"Hello World\",\"testInt\":128,\"testStr\":\"hello\",\"time\":\"logtimeplaceholder\"}\n",
			},
		},
		{
			name: "Special Level",
			args: logging.Formatter{
				LevelKey: "bubbleScale",
			},
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "bubbleScale=info msg=\"Hello World\" testInt=128 testStr=hello time=\"logtimeplaceholder\"\n",
				logging.JsonFormat: "{\"bubbleScale\":\"info\",\"msg\":\"Hello World\",\"testInt\":128,\"testStr\":\"hello\",\"time\":\"logtimeplaceholder\"}\n",
			},
		},
		{
			name: "Level Key",
			args: logging.Formatter{
				LevelKey: "~",
			},
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "level=info msg=\"Hello World\" testInt=128 testStr=hello time=\"logtimeplaceholder\"\n",
				logging.JsonFormat: "{\"level\":\"info\",\"msg\":\"Hello World\",\"testInt\":128,\"testStr\":\"hello\",\"time\":\"logtimeplaceholder\"}\n",
			},
		},
		{
			name: "Special Function",
			args: logging.Formatter{
				FunctionKey: "pickle",
			},
			postLogFunc: GetFileAndFunctionPlaceholder("pickle"),
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "file=\"fileplaceholder\" level=info msg=\"Hello World\" pickle=funcplaceholder testInt=128 testStr=hello time=\"logtimeplaceholder\"\n",
				logging.JsonFormat: "{\"file\":\"fileplaceholder\",\"level\":\"info\",\"msg\":\"Hello World\",\"pickle\":\"funcplaceholder\",\"testInt\":128,\"testStr\":\"hello\",\"time\":\"logtimeplaceholder\"}\n",
			},
		},
		{
			name: "Function Key",
			args: logging.Formatter{
				FunctionKey: "~",
			},
			postLogFunc: GetFileAndFunctionPlaceholder("function"),
			expected: map[logging.FormatterType]string{
				logging.TextFormat: "file=\"fileplaceholder\" function=funcplaceholder level=info msg=\"Hello World\" testInt=128 testStr=hello time=\"logtimeplaceholder\"\n",
				logging.JsonFormat: "{\"file\":\"fileplaceholder\",\"function\":\"funcplaceholder\",\"level\":\"info\",\"msg\":\"Hello World\",\"testInt\":128,\"testStr\":\"hello\",\"time\":\"logtimeplaceholder\"}\n",
			},
		},
	}
	for formatType := logging.TextFormat; formatType <= logging.JsonFormat; formatType++ {
		t.Run(fmt.Sprintf("Format-%s", formatType), func(t *testing.T) {
			for _, testFormat := range tests {
				t.Run(testFormat.name, func(t *testing.T) {
					logArgs := new(logging.TLMLoggingInitialization)
					logArgs.Formatter = logging.Formatter{
						Type:        formatType,
						TimeKey:     testFormat.args.TimeKey,
						MessageKey:  testFormat.args.MessageKey,
						LevelKey:    testFormat.args.LevelKey,
						FunctionKey: testFormat.args.FunctionKey,
						TimeFormat:  testFormat.args.TimeFormat,
					}
					logger, buffer := createLogger(logArgs)

					logTime := time.Now()
					logger.WithField("testStr", "hello").WithField("testInt", 128).Info("Hello World")

					expected, ok := testFormat.expected[formatType]
					util.AssertEqual(t, ok, true, "expected log output")

					newLog, parsedLogTime := SplitAndOrderOutput(buffer.String(), logArgs.Formatter)
					if testFormat.postLogFunc != nil {
						newLog = testFormat.postLogFunc(newLog)
					}
					AssertTime(t, logArgs.Formatter.TimeKey, logArgs.Formatter.TimeFormat, parsedLogTime, newLog, logTime)
					util.AssertEqual(t, newLog, expected, "contents")
				})
			}
		})
	}
}

func TestSanity(t *testing.T) {
	// This exists as a sanity check for some constants

	logArgs := new(logging.TLMLoggingInitialization)
	logger, _ := createLogger(logArgs)

	// For printing out the caller of a log, various packages do dynamic calculation. This is fine for many but times time and resources and best done in a Once or on init
	// As the result is always the same (unless you fork the repo), use a const instead. If you fork, this test will fail and you need to update LoggingPackageName
	// If this becomes big enough and has many forks, dynamic parsing can be looked into
	lrus, ok := logger.(*logging.LogrusImpl)
	util.AssertEqual(t, ok, true, "LogrusImpl")
	util.AssertEqual(t, reflect.TypeOf(*lrus).PkgPath(), logging.LoggingPackageName, "package name")

	// We're testing Logrus... if it's something other then Logrus, then this will fail
	util.AssertEqual(t, reflect.TypeOf(*lrus.Log).PkgPath(), logging.LogrusPackageName, "logrus package name")
}
