package logging_test

import (
	"errors"
	"testing"

	"github.com/rcmaniac25/tlm/logging"
	"github.com/rcmaniac25/tlm/util"
)

func TestInitLogging(t *testing.T) {
	type args struct {
		logArgs *logging.TLMLoggingInitialization
		setup   func()
		cleanup func()
	}
	type expected struct {
		expectError  bool
		expectLogger bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "No Args",
			args: args{
				logArgs: nil,
			},
			expected: expected{
				expectError:  false,
				expectLogger: false,
			},
		},
		{
			name: "Default Args",
			args: args{
				logArgs: &logging.TLMLoggingInitialization{},
			},
			expected: expected{
				expectError:  true,
				expectLogger: false,
			},
		},
		{
			name: "Invalid Type",
			args: args{
				logArgs: &logging.TLMLoggingInitialization{
					Type: -1,
				},
			},
			expected: expected{
				expectError:  true,
				expectLogger: false,
			},
		},
		{
			name: "Logrus",
			args: args{
				logArgs: &logging.TLMLoggingInitialization{
					Type: logging.LogrusLogType,
				},
			},
			expected: expected{
				expectError:  false,
				expectLogger: true,
			},
		},
		{
			name: "No Custom type",
			args: args{
				logArgs: &logging.TLMLoggingInitialization{
					Type: logging.CustomLogType,
				},
			},
			expected: expected{
				expectError:  true,
				expectLogger: false,
			},
		},
		{
			name: "Custom type",
			args: args{
				logArgs: &logging.TLMLoggingInitialization{
					Type:        logging.CustomLogType,
					CustomeType: "MyFakeLogger",
				},
				setup: func() {
					logging.RegisterLogger("MyFakeLogger", func(_ *logging.TLMLoggingInitialization) (logging.Logger, error) {
						return &logging.NullLogger, nil
					})
				},
				cleanup: func() {
					logging.UnregisterLogger("MyFakeLogger")
				},
			},
			expected: expected{
				expectError:  false,
				expectLogger: true,
			},
		},
		{
			name: "Custom type (failed)",
			args: args{
				logArgs: &logging.TLMLoggingInitialization{
					Type:        logging.CustomLogType,
					CustomeType: "MyFakeLogger",
				},
				setup: func() {
					logging.RegisterLogger("MyFakeLogger", func(_ *logging.TLMLoggingInitialization) (logging.Logger, error) {
						return nil, errors.New("ohnoes")
					})
				},
				cleanup: func() {
					logging.UnregisterLogger("MyFakeLogger")
				},
			},
			expected: expected{
				expectError:  true,
				expectLogger: false,
			},
		},
		{
			name: "Custom type (no logger)",
			args: args{
				logArgs: &logging.TLMLoggingInitialization{
					Type:        logging.CustomLogType,
					CustomeType: "MyFakeLogger",
				},
				setup: func() {
					logging.RegisterLogger("MyFakeLogger", func(_ *logging.TLMLoggingInitialization) (logging.Logger, error) {
						return nil, nil
					})
				},
				cleanup: func() {
					logging.UnregisterLogger("MyFakeLogger")
				},
			},
			expected: expected{
				expectError:  true,
				expectLogger: false,
			},
		},
		{
			name: "Custom type (wrong type)",
			args: args{
				logArgs: &logging.TLMLoggingInitialization{
					Type:        logging.CustomLogType,
					CustomeType: "MyFakeLogger2",
				},
				setup: func() {
					logging.RegisterLogger("MyFakeLogger", func(_ *logging.TLMLoggingInitialization) (logging.Logger, error) {
						return &logging.NullLogger, nil
					})
				},
				cleanup: func() {
					logging.UnregisterLogger("MyFakeLogger")
				},
			},
			expected: expected{
				expectError:  true,
				expectLogger: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args.setup != nil {
				test.args.setup()
			}
			if test.args.cleanup != nil {
				defer test.args.cleanup()
			}

			logger, err := logging.InitLogging(test.args.logArgs)
			if test.expected.expectError {
				util.AssertError(t, err, "error")
			} else {
				util.AssertNoError(t, err, "no error")
			}

			if test.expected.expectLogger {
				util.AssertNotEqual(t, logger, nil, "logger")
			} else {
				util.AssertEqual(t, logger, nil, "no logger")
			}
		})
	}
}

func TestRegisterLogger(t *testing.T) {
	type args struct {
		loggerName string
		initFunc   logging.CustomeLoggerInitializationFunc
		setup      func()
		cleanup    func()
	}
	tests := []struct {
		name          string
		args          args
		expectedError bool
	}{
		{
			name: "No Name",
			args: args{
				loggerName: "",
			},
			expectedError: true,
		},
		{
			name: "No Init Func",
			args: args{
				loggerName: "specialLogger",
			},
			expectedError: true,
		},
		{
			name: "Already Exists",
			args: args{
				loggerName: "specialLogger",
				initFunc: func(_ *logging.TLMLoggingInitialization) (logging.Logger, error) {
					return nil, nil
				},
				setup: func() {
					logging.RegisterLogger("specialLogger", func(_ *logging.TLMLoggingInitialization) (logging.Logger, error) {
						return nil, nil
					})
				},
				cleanup: func() {
					logging.UnregisterLogger("specialLogger")
				},
			},
			expectedError: true,
		},
		{
			name: "Registered",
			args: args{
				loggerName: "specialLogger",
				initFunc: func(_ *logging.TLMLoggingInitialization) (logging.Logger, error) {
					return nil, nil
				},
			},
			expectedError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer logging.UnregisterLogger(test.args.loggerName)

			if test.args.setup != nil {
				test.args.setup()
			}
			if test.args.cleanup != nil {
				defer test.args.cleanup()
			}

			err := logging.RegisterLogger(test.args.loggerName, test.args.initFunc)
			if test.expectedError {
				util.AssertError(t, err, "error")
			} else {
				util.AssertNoError(t, err, "no error")
			}
		})
	}
}

func TestUnregisterLogger(t *testing.T) {
	type args struct {
		loggerName string
		setup      func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Logger Doesn't Exist",
			args: args{
				loggerName: "iCease",
			},
		},
		{
			name: "Logger Does Exist",
			args: args{
				loggerName: "iLive",
				setup: func() {
					//TODO
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args.setup != nil {
				test.args.setup()
			}

			util.AssertNoPanic(t, func() {
				logging.UnregisterLogger(test.args.loggerName)
			}, "unregister")
		})
	}
}
