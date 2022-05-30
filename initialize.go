package tlm

import (
	"context"
	"errors"

	
)

func Startup(init *TLMInitialization) (context.Context, error) {
	if init == nil {
		return context.Background(), errors.New("initialization values cannot be nil")
	}

	logger, err := InitLogging(init.Logging)
	if err != nil {
		return context.Background(), err
	}

	//TODO
	return context.Background(), nil
}
