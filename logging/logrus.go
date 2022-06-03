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

func (r *LogrusImpl) Debugf(format string, args ...any) {
	r.LR.Debugf(format, args...)
}
func (r *LogrusImpl) Debug(args ...any) {
	r.LR.Debug(args...)
}
func (r *LogrusImpl) Debugln(args ...any) {
	r.LR.Debugln(args...)
}

func (r *LogrusImpl) Infof(format string, args ...any) {
	r.LR.Infof(format, args...)
}
func (r *LogrusImpl) Info(args ...any) {
	r.LR.Info(args...)
}
func (r *LogrusImpl) Infoln(args ...any) {
	r.LR.Infoln(args...)
}

func (r *LogrusImpl) Warnf(format string, args ...any) {
	r.LR.Warnf(format, args...)
}
func (r *LogrusImpl) Warn(args ...any) {
	r.LR.Warn(args...)
}
func (r *LogrusImpl) Warnln(args ...any) {
	r.LR.Warnln(args...)
}

func (r *LogrusImpl) Errorf(format string, args ...any) {
	r.LR.Errorf(format, args...)
}
func (r *LogrusImpl) Error(args ...any) {
	r.LR.Error(args...)
}
func (r *LogrusImpl) Errorln(args ...any) {
	r.LR.Errorln(args...)
}

func (r *LogrusImpl) Panicf(format string, args ...any) {
	r.LR.Panicf(format, args...)
}
func (r *LogrusImpl) Panic(args ...any) {
	r.LR.Panic(args...)
}
func (r *LogrusImpl) Panicln(args ...any) {
	r.LR.Panicln(args...)
}

func (r *LogrusImpl) Fatalf(format string, args ...any) {
	r.LR.Fatalf(format, args...)
}
func (r *LogrusImpl) Fatal(args ...any) {
	r.LR.Fatal(args...)
}
func (r *LogrusImpl) Fatalln(args ...any) {
	r.LR.Fatalln(args...)
}
