package logging

import (
	"fmt"
)

func InitLogging(args *TLMLoggingInitialization) (*Logger, error) {
	if args == nil {
		// If no logging info, then no need to process
		return nil, nil
	}

	log := new(Logger)

	var err error
	switch args.Type {
	case Logrus:
		log.Impl, err = InitLogrus(args)
	default:
		return nil, fmt.Errorf("unknown logging type: %v", args.Type)
	}
	
	if err != nil {
		return nil, err
	}

	//TODO

	return log, nil
}