package log

import (
	"bytes"
	"context"
	"fmt"
	"runtime"

	"github.com/rs/zerolog"
)

const (
	LEVEL_DEBUG = uint8(iota)
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR

	ZEROLOG_CALLER_SKIP = 3
	YKLOG_CALLER_SKIP   = 4
)

func init() {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.MessageFieldName = "msg"
}

type Fields = map[string]interface{}

func newLoggerContext() *LoggerContext {
	logger := &LoggerContext{}
	logger.loggerConf = LoggerConf{}
	return logger
}

type LoggerConf struct {
	LogFilePath string // Log文件全路径
	Level       uint8
	Caller      bool // 是否需要显示Caller，默认false时不显示（注：获取caller会影响性能）
	HideTime    bool // 是否隐藏时间，默认显示
	ContextCall func(l *LoggerContext) LoggerContextIface
}

type LoggerContextIface interface {
	WithCommonFields(commonFields Fields) LoggerContextIface
	Tag(tag string) LoggerContextIface

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	GetDebugLog() *zerolog.Event
	GetInfoLog() *zerolog.Event
	GetWarnLog() *zerolog.Event
	GetErrorLog() *zerolog.Event
}

type LoggerContext struct {
	logger *zerolog.Logger

	tag          string
	requestId    string
	commonFields []Fields
	fields       []Fields

	loggerConf LoggerConf

	ctx context.Context
}

func (l *LoggerContext) WithCommonFields(commonFields Fields) LoggerContextIface {
	logger := l.clone()
	logger.commonFields = append(logger.commonFields, commonFields)
	return logger
}

func (l *LoggerContext) Tag(tag string) LoggerContextIface {
	logger := l.clone()
	logger.tag = tag
	return logger
}

func (l *LoggerContext) Debug(args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprint(&b, args...)
	l.log(LEVEL_DEBUG, b.String())
}

func (l *LoggerContext) Debugf(format string, args ...interface{}) {
	l.log(LEVEL_DEBUG, fmt.Sprintf(format, args...))
}

func (l *LoggerContext) Info(args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprint(&b, args...)
	l.log(LEVEL_INFO, b.String())
}

func (l *LoggerContext) Infof(format string, args ...interface{}) {
	l.log(LEVEL_INFO, fmt.Sprintf(format, args...))
}

func (l *LoggerContext) Warn(args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprint(&b, args...)
	l.log(LEVEL_WARN, b.String())
}

func (l *LoggerContext) Warnf(format string, args ...interface{}) {
	l.log(LEVEL_WARN, fmt.Sprintf(format, args...))
}

func (l *LoggerContext) Error(args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprint(&b, args...)
	l.log(LEVEL_ERROR, b.String())
}

func (l *LoggerContext) Errorf(format string, args ...interface{}) {
	l.log(LEVEL_ERROR, fmt.Sprintf(format, args...))
}

func (l *LoggerContext) GetDebugLog() *zerolog.Event {
	loggerEvent := l.logger.Debug()
	l.withCommon(loggerEvent, ZEROLOG_CALLER_SKIP)
	return loggerEvent
}
func (l *LoggerContext) GetInfoLog() *zerolog.Event {
	loggerEvent := l.logger.Info()
	l.withCommon(loggerEvent, ZEROLOG_CALLER_SKIP)
	return loggerEvent
}
func (l *LoggerContext) GetWarnLog() *zerolog.Event {
	loggerEvent := l.logger.Warn()
	l.withCommon(loggerEvent, ZEROLOG_CALLER_SKIP)
	return loggerEvent
}
func (l *LoggerContext) GetErrorLog() *zerolog.Event {
	loggerEvent := l.logger.Error()
	l.withCommon(loggerEvent, ZEROLOG_CALLER_SKIP)
	return loggerEvent
}

func (l *LoggerContext) log(level uint8, msg string) {
	var zeroEvent *zerolog.Event
	switch level {
	case LEVEL_DEBUG:
		zeroEvent = l.logger.Debug()
	case LEVEL_INFO:
		zeroEvent = l.logger.Info()
	case LEVEL_WARN:
		zeroEvent = l.logger.Warn()
	case LEVEL_ERROR:
		zeroEvent = l.logger.Error()
	}
	l.withCommon(zeroEvent, YKLOG_CALLER_SKIP)
	zeroEvent.Msg(msg)
}

func (l *LoggerContext) generateCallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d %s", file, line, f.Name())
}

func (l *LoggerContext) clone() *LoggerContext {
	newLog := *l
	return &newLog
}

func (l *LoggerContext) withCommon(loggerEvent *zerolog.Event, skip int) {
	for _, field := range l.commonFields {
		loggerEvent.Fields(field)
	}

	if !l.loggerConf.HideTime {
		loggerEvent.Timestamp()
	}
	if l.tag != "" {
		loggerEvent.Str("tag", l.tag)
	}
	if l.requestId != "" {
		loggerEvent.Str("request-id", l.requestId)
	}
	if l.loggerConf.Caller {
		loggerEvent.Str("caller", l.generateCallerInfo(skip))
	}
	dict := zerolog.Dict()
	for _, field := range l.fields {
		dict.Fields(field)
	}
	if len(l.fields) > 0 {
		loggerEvent.Dict("attach", dict)
	}
}
