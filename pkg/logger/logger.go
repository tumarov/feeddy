package logger

type Logger interface {
	InvalidArg(name string)
	InvalidArgValue(name string, value string)
	MissingArgMessage(name string)
	Exception(err error)
	Debugf(format string, args ...interface{})
}
