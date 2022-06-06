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

	if level, ok := convertLogLevel(args.Level); ok {
		logger.LR.SetLevel(level)
	}

	if formatter, ok := getFormatter(args.Formatter, logger.LR.Formatter); ok {
		logger.LR.Formatter = formatter
	}

	//TODO

	return logger, nil
}

func convertLogLevel(level LogLevel) (logrus.Level, bool) {
	switch level {
	case DebugLevel:
		return logrus.DebugLevel, true
	case InfoLevel:
		return logrus.InfoLevel, true
	case WarnLevel:
		return logrus.WarnLevel, true
	case ErrorLevel:
		return logrus.ErrorLevel, true
	case PanicLevel:
		return logrus.PanicLevel, true
	case FatalLevel:
		return logrus.FatalLevel, true
	}
	return logrus.InfoLevel, false
}

// Formatting
func getFormatter(formatterArgs Formatter, def logrus.Formatter) (logrus.Formatter, bool) {
	switch formatterArgs.Type {
	case TextFormat:
		return getTextFormatter(formatterArgs, nil), true
	case JsonFormat:
		return getJsonFormatter(formatterArgs, nil), true
	case DefaultFormat:
		if text, ok := def.(*logrus.TextFormatter); ok {
			return getTextFormatter(formatterArgs, text), true
		}
		if json, ok := def.(*logrus.JSONFormatter); ok {
			return getJsonFormatter(formatterArgs, json), true
		}
	}
	return nil, false
}

func getTextFormatter(formatterArgs Formatter, form *logrus.TextFormatter) logrus.Formatter {
	if form == nil {
		form = new(logrus.TextFormatter)
	}
	form.FieldMap = setFormatterFieldMap(formatterArgs)
	if formatterArgs.TimeKey == "-" {
		form.DisableTimestamp = true
	}
	form.TimestampFormat = formatterArgs.TimeFormat
	return form
}

func getJsonFormatter(formatterArgs Formatter, form *logrus.JSONFormatter) logrus.Formatter {
	if form == nil {
		form = new(logrus.JSONFormatter)
	}
	form.FieldMap = setFormatterFieldMap(formatterArgs)
	if formatterArgs.TimeKey == "-" {
		form.DisableTimestamp = true
	}
	form.TimestampFormat = formatterArgs.TimeFormat
	return form
}

func setFormatterFieldMap(formatterArgs Formatter) logrus.FieldMap {
	fieldMap := make(logrus.FieldMap)

	switch formatterArgs.TimeKey {
	// "-" will skip the time and is handled outside of this function
	case "~":
		fieldMap[logrus.FieldKeyTime] = "Time"
	case "":
	default:
		fieldMap[logrus.FieldKeyTime] = formatterArgs.TimeKey
	}

	switch formatterArgs.MessageKey {
	case "~":
		fieldMap[logrus.FieldKeyMsg] = "Message"
	case "":
	default:
		fieldMap[logrus.FieldKeyMsg] = formatterArgs.MessageKey
	}

	switch formatterArgs.LevelKey {
	case "~":
		fieldMap[logrus.FieldKeyLevel] = "Level"
	case "":
	default:
		fieldMap[logrus.FieldKeyLevel] = formatterArgs.LevelKey
	}

	switch formatterArgs.FunctionKey {
	case "~":
		fieldMap[logrus.FieldKeyFunc] = "Function"
	case "":
	default:
		fieldMap[logrus.FieldKeyFunc] = formatterArgs.FunctionKey
	}

	return fieldMap
}

// Logging function calls
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
