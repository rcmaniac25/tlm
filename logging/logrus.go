package logging

import (
	"github.com/rcmaniac25/tlm/util"

	"github.com/sirupsen/logrus"
)

type LogrusImpl struct {
	Log   *logrus.Logger
	Entry *logrus.Entry
}

func InitLogrus(args *TLMLoggingInitialization) (Logger, error) {
	logger := &LogrusImpl{
		Log: logrus.New(),
	}

	if args.Output != nil {
		logger.Log.SetOutput(args.Output)
	}

	if level, ok := convertLogLevel(args.Level); ok {
		logger.Log.SetLevel(level)
	}

	if formatter, ok := getFormatter(args.Formatter, logger.Log.Formatter); ok {
		logger.Log.Formatter = formatter
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
		return getTextFormatter(formatterArgs, nil)
	case JsonFormat:
		return getJsonFormatter(formatterArgs, nil)
	case DefaultFormat:
		if text, ok := def.(*logrus.TextFormatter); ok {
			return getTextFormatter(formatterArgs, text)
		}
		if json, ok := def.(*logrus.JSONFormatter); ok {
			return getJsonFormatter(formatterArgs, json)
		}
	}
	return nil, false
}

func getTextFormatter(formatterArgs Formatter, form *logrus.TextFormatter) (logrus.Formatter, bool) {
	dirty := false
	if form == nil {
		form = new(logrus.TextFormatter)
		dirty = true
	}
	fields, fieldsDirty := setFormatterFieldMap(formatterArgs)
	if fieldsDirty {
		form.FieldMap = fields
		dirty = true
	}
	if formatterArgs.TimeKey == "-" && !form.DisableTimestamp {
		form.DisableTimestamp = true
		dirty = true
	}
	dirty = dirty || (form.TimestampFormat != formatterArgs.TimeFormat)
	form.TimestampFormat = formatterArgs.TimeFormat
	return form, dirty
}

func getJsonFormatter(formatterArgs Formatter, form *logrus.JSONFormatter) (logrus.Formatter, bool) {
	dirty := false
	if form == nil {
		form = new(logrus.JSONFormatter)
		dirty = true
	}
	fields, fieldsDirty := setFormatterFieldMap(formatterArgs)
	if fieldsDirty {
		form.FieldMap = fields
		dirty = true
	}
	if formatterArgs.TimeKey == "-" && !form.DisableTimestamp {
		form.DisableTimestamp = true
		dirty = true
	}
	dirty = dirty || (form.TimestampFormat != formatterArgs.TimeFormat)
	form.TimestampFormat = formatterArgs.TimeFormat
	return form, dirty
}

func setFormatterFieldMap(formatterArgs Formatter) (logrus.FieldMap, bool) {
	fieldMap := make(logrus.FieldMap)

	dirty := formatterArgs.TimeKey != ""
	switch formatterArgs.TimeKey {
	// "-" will skip the time and is handled outside of this function
	case "~":
		fieldMap[logrus.FieldKeyTime] = "Time"
	case "":
	default:
		fieldMap[logrus.FieldKeyTime] = formatterArgs.TimeKey
	}

	dirty = dirty || formatterArgs.MessageKey != ""
	switch formatterArgs.MessageKey {
	case "~":
		fieldMap[logrus.FieldKeyMsg] = "Message"
	case "":
	default:
		fieldMap[logrus.FieldKeyMsg] = formatterArgs.MessageKey
	}

	dirty = dirty || formatterArgs.LevelKey != ""
	switch formatterArgs.LevelKey {
	case "~":
		fieldMap[logrus.FieldKeyLevel] = "Level"
	case "":
	default:
		fieldMap[logrus.FieldKeyLevel] = formatterArgs.LevelKey
	}

	dirty = dirty || formatterArgs.FunctionKey != ""
	switch formatterArgs.FunctionKey {
	case "~":
		fieldMap[logrus.FieldKeyFunc] = "Function"
	case "":
	default:
		fieldMap[logrus.FieldKeyFunc] = formatterArgs.FunctionKey
	}

	return fieldMap, dirty
}

// Fields
func (r *LogrusImpl) WithField(key string, value any) Logger {
	if r.Entry != nil {
		return &LogrusImpl{Entry: r.Entry.WithField(key, value)}
	}
	return &LogrusImpl{Entry: r.Log.WithField(key, value)}
}

func (r *LogrusImpl) WithFields(fields util.Fields) Logger {
	logFields := make(logrus.Fields)
	for key, value := range fields {
		logFields[key] = value
	}
	if r.Entry != nil {
		return &LogrusImpl{Entry: r.Entry.WithFields(logFields)}
	}
	return &LogrusImpl{Entry: r.Log.WithFields(logFields)}
}

// Logging function calls
func (r *LogrusImpl) Debugf(format string, args ...any) {
	if r.Entry != nil {
		r.Entry.Debugf(format, args...)
		return
	}
	r.Log.Debugf(format, args...)
}
func (r *LogrusImpl) Debug(args ...any) {
	if r.Entry != nil {
		r.Entry.Debug(args...)
		return
	}
	r.Log.Debug(args...)
}
func (r *LogrusImpl) Debugln(args ...any) {
	if r.Entry != nil {
		r.Entry.Debugln(args...)
		return
	}
	r.Log.Debugln(args...)
}

func (r *LogrusImpl) Infof(format string, args ...any) {
	if r.Entry != nil {
		r.Entry.Infof(format, args...)
		return
	}
	r.Log.Infof(format, args...)
}
func (r *LogrusImpl) Info(args ...any) {
	if r.Entry != nil {
		r.Entry.Info(args...)
		return
	}
	r.Log.Info(args...)
}
func (r *LogrusImpl) Infoln(args ...any) {
	if r.Entry != nil {
		r.Entry.Infoln(args...)
		return
	}
	r.Log.Infoln(args...)
}

func (r *LogrusImpl) Warnf(format string, args ...any) {
	if r.Entry != nil {
		r.Entry.Warnf(format, args...)
		return
	}
	r.Log.Warnf(format, args...)
}
func (r *LogrusImpl) Warn(args ...any) {
	if r.Entry != nil {
		r.Entry.Warn(args...)
		return
	}
	r.Log.Warn(args...)
}
func (r *LogrusImpl) Warnln(args ...any) {
	if r.Entry != nil {
		r.Entry.Warnln(args...)
		return
	}
	r.Log.Warnln(args...)
}

func (r *LogrusImpl) Errorf(format string, args ...any) {
	if r.Entry != nil {
		r.Entry.Errorf(format, args...)
		return
	}
	r.Log.Errorf(format, args...)
}
func (r *LogrusImpl) Error(args ...any) {
	if r.Entry != nil {
		r.Entry.Error(args...)
		return
	}
	r.Log.Error(args...)
}
func (r *LogrusImpl) Errorln(args ...any) {
	if r.Entry != nil {
		r.Entry.Errorln(args...)
		return
	}
	r.Log.Errorln(args...)
}

func (r *LogrusImpl) Panicf(format string, args ...any) {
	if r.Entry != nil {
		r.Entry.Panicf(format, args...)
		return
	}
	r.Log.Panicf(format, args...)
}
func (r *LogrusImpl) Panic(args ...any) {
	if r.Entry != nil {
		r.Entry.Panic(args...)
		return
	}
	r.Log.Panic(args...)
}
func (r *LogrusImpl) Panicln(args ...any) {
	if r.Entry != nil {
		r.Entry.Panicln(args...)
		return
	}
	r.Log.Panicln(args...)
}

func (r *LogrusImpl) Fatalf(format string, args ...any) {
	if r.Entry != nil {
		r.Entry.Fatalf(format, args...)
		return
	}
	r.Log.Fatalf(format, args...)
}
func (r *LogrusImpl) Fatal(args ...any) {
	if r.Entry != nil {
		r.Entry.Fatal(args...)
		return
	}
	r.Log.Fatal(args...)
}
func (r *LogrusImpl) Fatalln(args ...any) {
	if r.Entry != nil {
		r.Entry.Fatalln(args...)
		return
	}
	r.Log.Fatalln(args...)
}
