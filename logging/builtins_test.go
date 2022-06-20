package logging_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rcmaniac25/tlm"
	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
)

type BuiltinLogger struct {
	Name             string
	Logger           logging.Logger
	Collector        *logging.DebugLogCollector
	FatalLogsHandled bool
}

func getLoggers() []BuiltinLogger {
	loggers := []struct {
		name         string
		logGenerated func() (logging.Logger, *logging.DebugLogCollector, bool)
		logMapper    func(logging.LogLevel) string
	}{
		{
			name: "Null Logger",
			logGenerated: func() (logging.Logger, *logging.DebugLogCollector, bool) {
				return tlm.Log(context.Background()), nil, false
			},
		},
		{
			name: "Logrus",
			logGenerated: func() (logging.Logger, *logging.DebugLogCollector, bool) {
				inits := new(tlm.TLMInitialization)
				inits.Logging = new(logging.TLMLoggingInitialization)

				collector := logging.NewDebugLogCollector()
				collector.SetupInitialization(inits.Logging)
				inits.Logging.Level = logging.DebugLevel

				ctx, err := tlm.Startup(inits)
				if err != nil {
					return nil, nil, false
				}

				logger := tlm.Log(ctx)

				exitHandlerSet := false
				type InternalTestingExitHandler interface {
					TestingSetFatalExitFunction(exitHandler func(int)) bool
				}
				if v, ok := logger.(InternalTestingExitHandler); ok {
					exitHandlerSet = v.TestingSetFatalExitFunction(collector.OnExitCode)
				}

				return logger, collector, exitHandlerSet
			},
			logMapper: func(tlmLogLevel logging.LogLevel) string {
				if tlmLogLevel == logging.WarnLevel {
					return "warning"
				}
				return tlmLogLevel.String()
			},
		},
	}
	results := make([]BuiltinLogger, 0)
	for _, logger := range loggers {
		log, collector, exitHandlerSet := logger.logGenerated()
		results = append(results, BuiltinLogger{
			Name:             logger.name,
			Logger:           log,
			Collector:        collector,
			FatalLogsHandled: exitHandlerSet,
		})
	}
	return results
}

func TestWithField(t *testing.T) {
	for _, logItem := range getLoggers() {
		t.Run(logItem.Name, func(t *testing.T) {
			util.AssertNoPanic(t, func() {
				logItem.Logger.WithField("myField", 1234).Info("Test")
				logItem.Logger.WithField("myField", 1234).WithField("theOtherOne", "soup").Info("Test2")
			}, "log")
			if logItem.Collector != nil {
				util.AssertEqual(t, logItem.Collector.GetNumberLogs(), 2, "count")

				util.AssertEqual(t, logItem.Collector.GetMessage(0), "Test", "message")
				util.AssertEqualExistsFunc(t, logItem.Collector.GetFieldFunc(0, "myField"), float64(1234), "field")

				util.AssertEqual(t, logItem.Collector.GetMessage(1), "Test2", "message")
				util.AssertEqualExistsFunc(t, logItem.Collector.GetFieldFunc(1, "myField"), float64(1234), "field")
				util.AssertEqualExistsFunc(t, logItem.Collector.GetFieldFunc(1, "theOtherOne"), "soup", "field")
			}
		})
	}
}

func TestWithFields(t *testing.T) {
	for _, logItem := range getLoggers() {
		t.Run(logItem.Name, func(t *testing.T) {
			util.AssertNoPanic(t, func() {
				logItem.Logger.WithFields(util.Fields{
					"hi":      "people",
					"myField": 6547,
				}).Info("Test")

				logItem.Logger.WithFields(util.Fields{
					"hi":      "people",
					"myField": 6547,
				}).WithFields(util.Fields{
					"reason": errors.New("Ohai there"),
				}).Info("Test")
			}, "log")
			if logItem.Collector != nil {
				util.AssertEqual(t, logItem.Collector.GetNumberLogs(), 2, "count")

				util.AssertEqual(t, logItem.Collector.GetMessage(0), "Test", "message")
				util.AssertEqualExistsFunc(t, logItem.Collector.GetFieldFunc(0, "myField"), float64(6547), "field")
				util.AssertEqualExistsFunc(t, logItem.Collector.GetFieldFunc(0, "hi"), "people", "field")

				util.AssertEqual(t, logItem.Collector.GetMessage(1), "Test", "message")
				util.AssertEqualExistsFunc(t, logItem.Collector.GetFieldFunc(1, "myField"), float64(6547), "field")
				util.AssertEqualExistsFunc(t, logItem.Collector.GetFieldFunc(1, "hi"), "people", "field")
				util.AssertEqualExistsFunc(t, logItem.Collector.GetFieldFunc(1, "reason"), "Ohai there", "field")
			}
		})
	}
}

func TestLoggingFunctions(t *testing.T) {
	type args struct {
		logFunc    func(logging.Logger)
		fieldValue string
		willPanic  bool
	}
	type expected struct {
		msg   string
		level logging.LogLevel
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "Debug",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Debug("DebugTest", 123, float64(10.123), "Value") },
				fieldValue: "DebugField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "DebugTest123 10.123Value",
				level: logging.DebugLevel,
			},
		},
		{
			name: "Debugf",
			args: args{
				logFunc: func(logger logging.Logger) {
					logger.Debugf("%s-%d--%.3f---%s", "DebugTest", 123, float64(10.123), "Value")
				},
				fieldValue: "DebugField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "DebugTest-123--10.123---Value",
				level: logging.DebugLevel,
			},
		},
		{
			name: "Debugln",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Debugln("DebugTest", 123, float64(10.123), "Value") },
				fieldValue: "DebugField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "DebugTest 123 10.123 Value",
				level: logging.DebugLevel,
			},
		},
		{
			name: "Info",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Info("InfoTest", 123, float64(10.123), "Value") },
				fieldValue: "InfoField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "InfoTest123 10.123Value",
				level: logging.InfoLevel,
			},
		},
		{
			name: "Infof",
			args: args{
				logFunc: func(logger logging.Logger) {
					logger.Infof("%s-%d--%.3f---%s", "InfoTest", 123, float64(10.123), "Value")
				},
				fieldValue: "InfoField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "InfoTest-123--10.123---Value",
				level: logging.InfoLevel,
			},
		},
		{
			name: "Infoln",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Infoln("InfoTest", 123, float64(10.123), "Value") },
				fieldValue: "InfoField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "InfoTest 123 10.123 Value",
				level: logging.InfoLevel,
			},
		},
		{
			name: "Warn",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Warn("WarnTest", 123, float64(10.123), "Value") },
				fieldValue: "WarnField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "WarnTest123 10.123Value",
				level: logging.WarnLevel,
			},
		},
		{
			name: "Warnf",
			args: args{
				logFunc: func(logger logging.Logger) {
					logger.Warnf("%s-%d--%.3f---%s", "WarnTest", 123, float64(10.123), "Value")
				},
				fieldValue: "WarnField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "WarnTest-123--10.123---Value",
				level: logging.WarnLevel,
			},
		},
		{
			name: "Warnln",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Warnln("WarnTest", 123, float64(10.123), "Value") },
				fieldValue: "WarnField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "WarnTest 123 10.123 Value",
				level: logging.WarnLevel,
			},
		},
		{
			name: "Error",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Error("ErrorTest", 123, float64(10.123), "Value") },
				fieldValue: "ErrorField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "ErrorTest123 10.123Value",
				level: logging.ErrorLevel,
			},
		},
		{
			name: "Errorf",
			args: args{
				logFunc: func(logger logging.Logger) {
					logger.Errorf("%s-%d--%.3f---%s", "ErrorTest", 123, float64(10.123), "Value")
				},
				fieldValue: "ErrorField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "ErrorTest-123--10.123---Value",
				level: logging.ErrorLevel,
			},
		},
		{
			name: "Errorln",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Errorln("ErrorTest", 123, float64(10.123), "Value") },
				fieldValue: "ErrorField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "ErrorTest 123 10.123 Value",
				level: logging.ErrorLevel,
			},
		},
		{
			name: "Panic",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Panic("PanicTest", 123, float64(10.123), "Value") },
				fieldValue: "PanicField",
				willPanic:  true,
			},
			expected: expected{
				msg:   "PanicTest123 10.123Value",
				level: logging.PanicLevel,
			},
		},
		{
			name: "Panicf",
			args: args{
				logFunc: func(logger logging.Logger) {
					logger.Panicf("%s-%d--%.3f---%s", "PanicTest", 123, float64(10.123), "Value")
				},
				fieldValue: "PanicField",
				willPanic:  true,
			},
			expected: expected{
				msg:   "PanicTest-123--10.123---Value",
				level: logging.PanicLevel,
			},
		},
		{
			name: "Panicln",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Panicln("PanicTest", 123, float64(10.123), "Value") },
				fieldValue: "PanicField",
				willPanic:  true,
			},
			expected: expected{
				msg:   "PanicTest 123 10.123 Value",
				level: logging.PanicLevel,
			},
		},
		{
			name: "Fatal",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Fatal("FatalTest", 123, float64(10.123), "Value") },
				fieldValue: "FatalField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "FatalTest123 10.123Value",
				level: logging.FatalLevel,
			},
		},
		{
			name: "Fatalf",
			args: args{
				logFunc: func(logger logging.Logger) {
					logger.Fatalf("%s-%d--%.3f---%s", "FatalTest", 123, float64(10.123), "Value")
				},
				fieldValue: "FatalField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "FatalTest-123--10.123---Value",
				level: logging.FatalLevel,
			},
		},
		{
			name: "Fatalln",
			args: args{
				logFunc:    func(logger logging.Logger) { logger.Fatalln("FatalTest", 123, float64(10.123), "Value") },
				fieldValue: "FatalField",
				willPanic:  false,
			},
			expected: expected{
				msg:   "FatalTest 123 10.123 Value",
				level: logging.FatalLevel,
			},
		},
	}
	for _, logger := range getLoggers() {
		t.Run(logger.Name, func(t *testing.T) {
			if !logger.FatalLogsHandled {
				t.Log("<<<<< Can't test fatal >>>>>")
			}
			for _, log := range tests {
				t.Run(log.name, func(t *testing.T) {
					if log.expected.level == logging.FatalLevel && !logger.FatalLogsHandled {
						t.SkipNow()
						return
					}

					defer func() {
						if logger.Collector != nil {
							logger.Collector.Clear()
						}
					}()

					// Logging
					runMsg := "panic"
					runner := util.AssertNoPanic
					if log.args.willPanic {
						runMsg = "no panic"
						runner = util.AssertPanic
					}
					runner(t, func() {
						log.args.logFunc(logger.Logger)
					}, runMsg)
					runner(t, func() {
						log.args.logFunc(logger.Logger.WithField("myField", log.args.fieldValue))
					}, runMsg)

					// Check
					if logger.Collector != nil {
						util.AssertEqual(t, logger.Collector.GetNumberLogs(), 2, "count")

						util.AssertEqual(t, logger.Collector.GetMessage(0), log.expected.msg, "msg 0")
						util.AssertEqual(t, logger.Collector.GetMessage(1), log.expected.msg, "msg 1")

						util.AssertEqual(t, logger.Collector.GetLogLevel(0), log.expected.level, "level 0")
						util.AssertEqual(t, logger.Collector.GetLogLevel(1), log.expected.level, "level 1")

						if log.expected.level == logging.FatalLevel && logger.FatalLogsHandled {
							util.AssertEqualExistsFunc(t, logger.Collector.GetFatalExitcodeFunc(0), 1, "fatal exit code 0")
							util.AssertEqualExistsFunc(t, logger.Collector.GetFatalExitcodeFunc(1), 1, "fatal exit code 1")
						}

						util.AssertEqualExistsFunc(t, logger.Collector.GetFieldFunc(1, "myField"), log.args.fieldValue, "field")
					}
				})
			}
		})
	}
}
