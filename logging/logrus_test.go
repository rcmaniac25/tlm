package logging_test

import (
	"bytes"
	"testing"

	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
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

//TODO: formatting (combo of tests-default, text, json / different "keys", time format)
