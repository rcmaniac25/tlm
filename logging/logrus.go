package logging

import "github.com/sirupsen/logrus"

type LogrusImpl struct {
	LR *logrus.Logger
}

func InitLogrus(args *TLMLoggingInitialization) (Logger, error) {
	var logger Logger
	logger = &LogrusImpl{
		LR: logrus.New(),
	}

	//TODO

	return logger, nil
}

func (r *LogrusImpl) Info(msg string) {
	r.LR.Info(msg)
}
