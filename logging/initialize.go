package logging

import (
	"errors"
	"fmt"
)

func InitLogging(args *TLMLoggingInitialization) (TLMLogger, error) {
	if args == nil {
		// If no logging info, then no need to process
		return nil, nil
	}

	//TODO: expectation is the logger is going to be a higher-level type and the implementations will be lower level
	//      this way some stuff (like WithField) can be handled without needing the implementations needing to do the work.
	//      They can just get the final "set this structured log value/field thing on whatever your log"

	var log Logger
	var err error

	switch args.Type {
	case Custom:
		if len(args.CustomeType) == 0 {
			return nil, errors.New("type 'Custom' requires 'CustomeType' to be set")
		}
		if customLoggerInit, ok := registeredLoggers[args.CustomeType]; ok {
			log, err = customLoggerInit(args)
		} else {
			err = fmt.Errorf("custom type could is not registered: %s", args.CustomeType)
		}
	case Logrus:
		log, err = InitLogrus(args)
	default:
		return nil, fmt.Errorf("unknown logging type: %v", args.Type)
	}
	if err != nil {
		return nil, err
	}

	return &selfReferentialLogger{
		LoggerImpl: log,
	}, nil
}

var registeredLoggers = make(map[string]CustomeLoggerInitializationFunc)

func RegisterLogger(typeName string, loggerInit CustomeLoggerInitializationFunc) error {
	if len(typeName) == 0 {
		return errors.New("typeName must be set")
	}
	if loggerInit == nil {
		return errors.New("loggerInit cannot be nil")
	}
	if _, ok := registeredLoggers[typeName]; ok {
		return fmt.Errorf("logger of type '%s' already registered", typeName)
	}
	registeredLoggers[typeName] = loggerInit
	return nil
}

func UnregisterLogger(typeName string) {
	delete(registeredLoggers, typeName)
}
