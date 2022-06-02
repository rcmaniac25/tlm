package logging

import "github.com/sirupsen/logrus"

type LogrusImpl struct {
	LR *logrus.Logger
}

func InitLogrus(args *TLMLoggingInitialization) (Logger, error) {
	logger := &LogrusImpl{
		LR: logrus.New(),
	}

	if args.Output != nil {
		logger.LR.SetOutput(args.Output)
	}

	//TODO

	return logger, nil
}

func (r *LogrusImpl) Info(msg string) {
	r.LR.Info(msg)
}
