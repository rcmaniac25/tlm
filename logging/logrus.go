package logging

import "github.com/sirupsen/logrus"

type LogrusImpl struct {
	Logger *logrus.Logger
}

func InitLogrus(args *TLMLoggingInitialization) (*LogrusImpl, error) {
	impl := new(LogrusImpl)
	impl.Logger = logrus.New()

	//TODO

	return impl, nil
}
