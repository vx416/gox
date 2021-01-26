package log

import "context"

type (
	// LoggerKey define context key for logger
	LoggerKey struct{}
)

var _logger Logger

func init() {
	_logger = newDefaultZap()
}

// Logger interface
type Logger interface {
	Err(err error) Logger
	Fields(fields map[string]interface{}) Logger
	Field(key string, val interface{}) Logger
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
	Caller(stack int) Logger

	Attach(ctx context.Context) context.Context
	Printf(msg string, args ...interface{})
}

// SetGlobal set global logger
func SetGlobal(logger Logger) Logger {
	_logger = logger
	return _logger
}

// Get get global logger
func Get() Logger {
	return _logger
}

// Attach attach logger instance into context
func Attach(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, LoggerKey{}, logger)
}

// Ctx get logger instacne from context
func Ctx(ctx context.Context) Logger {
	val := ctx.Value(LoggerKey{})
	logger, ok := val.(Logger)
	if !ok {
		return Get()
	}

	return logger
}
