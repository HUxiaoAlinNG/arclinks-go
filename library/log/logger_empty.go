package log

import (
	"fmt"

	"github.com/rs/zerolog"
)

type EmptyLoggerContext struct {
}

func (e *EmptyLoggerContext) WithCommonFields(fields Fields) LoggerContextIface {
	return e
}

func (e *EmptyLoggerContext) Tag(tag string) LoggerContextIface {
	return e
}

func (e *EmptyLoggerContext) Debug(msg ...interface{})                  {}
func (e *EmptyLoggerContext) Debugf(format string, args ...interface{}) {}
func (e *EmptyLoggerContext) Info(msg ...interface{})                   {}
func (e *EmptyLoggerContext) Infof(format string, args ...interface{})  {}
func (e *EmptyLoggerContext) Warn(msg ...interface{})                   {}
func (e *EmptyLoggerContext) Warnf(format string, args ...interface{})  {}
func (e *EmptyLoggerContext) Error(msg ...interface{}) {
	fmt.Print(msg)
}
func (e *EmptyLoggerContext) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (e *EmptyLoggerContext) GetDebugLog() *zerolog.Event {
	event := &zerolog.Event{}
	event.Discard()
	return event
}

func (e *EmptyLoggerContext) GetInfoLog() *zerolog.Event {
	event := &zerolog.Event{}
	event.Discard()
	return event
}

func (e *EmptyLoggerContext) GetWarnLog() *zerolog.Event {
	event := &zerolog.Event{}
	event.Discard()
	return event
}

func (e *EmptyLoggerContext) GetErrorLog() *zerolog.Event {
	event := &zerolog.Event{}
	event.Discard()
	return event
}
