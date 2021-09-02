package builtin

import (
	"github.com/sirupsen/logrus"
)

type Event struct {
	id      int
	message string
}

type Logger struct {
	*logrus.Logger
}

func NewBuiltinLogger() *Logger {
	baseLogger := logrus.New()
	builtinLogger := &Logger{baseLogger}
	builtinLogger.Formatter = &logrus.JSONFormatter{}
	return builtinLogger
}

var (
	invalidArgMessage      = Event{1, "invalid arg: %s"}
	invalidArgValueMessage = Event{2, "invalid value for argument: %s: %v"}
	missingArgMessage      = Event{3, "missing arg: %s"}
)

func (b *Logger) InvalidArg(name string) {
	b.Errorf(invalidArgMessage.message, name)
}

func (b *Logger) InvalidArgValue(name string, value string) {
	b.Errorf(invalidArgValueMessage.message, name, value)
}

func (b *Logger) MissingArgMessage(name string) {
	b.Errorf(missingArgMessage.message, name)
}

func (b *Logger) Exception(err error) {
	b.Error(err)
}

func (b *Logger) Debugf(format string, args ...interface{}) {
	b.Printf(format, args...)
}
