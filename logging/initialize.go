package logging

import (
	"fmt"
)

func InitLogging(args *TLMLoggingInitialization) (Logger, error) {
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
	case Logrus:
		log, err = InitLogrus(args)
	default:
		return nil, fmt.Errorf("unknown logging type: %v", args.Type)
	}

	if err != nil {
		return nil, err
	}

	//TODO

	return log, nil
}
