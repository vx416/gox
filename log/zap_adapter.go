package log

import (
	"context"
	"fmt"
	"time"

	"github.com/vx416/gox/log/encoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newDefaultZap() Logger {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zapConfig.EncoderConfig.EncodeLevel = encoder.ColorfulLevelEncoder
	zapConfig.EncoderConfig.EncodeCaller = encoder.ColorizeCallerEncoder
	logger, _ := zapConfig.Build()
	return &ZapAdapter{
		env:    "dev",
		zaplog: logger,
	}
}

type ZapAdapter struct {
	zaplog *zap.Logger
	env    Env
}

// Err implement logger interface
func (l *ZapAdapter) Err(err error) Logger {
	var (
		stackErr = GetStack(err)
		fields   = make([]zap.Field, 0, 2)
	)

	fields = append(fields, zap.String("error", err.Error()))
	if stackErr != nil {
		stackInfo := stackErr.String()
		fields = append(fields, zap.String("error_stack", stackInfo))
	}

	zaplog := l.zaplog.With(fields...)
	return l.clone(zaplog)
}

// Field implement logger interface
func (l *ZapAdapter) Field(key string, val interface{}) Logger {
	field := getField(key, val)
	zaplog := l.zaplog.With(field)
	return l.clone(zaplog)
}

// Fields implement logger interface
func (l *ZapAdapter) Fields(fieldsMap map[string]interface{}) Logger {
	fields := make([]zap.Field, 0, len(fieldsMap))
	for k, v := range fieldsMap {
		fields = append(fields, getField(k, v))
	}
	zaplog := l.zaplog.With(fields...)
	return l.clone(zaplog)
}

// Info implement logger interface
func (l *ZapAdapter) Info(msg string) {
	l.zaplog.Info(msg)
	defer l.zaplog.Sync()
}

// Debug implement logger interface
func (l *ZapAdapter) Debug(msg string) {
	l.zaplog.Debug(msg)
	defer l.zaplog.Sync()

}

// Warn implement logger interface
func (l *ZapAdapter) Warn(msg string) {
	l.zaplog.Warn(msg)
	defer l.zaplog.Sync()

}

// Error implement logger interface
func (l *ZapAdapter) Error(msg string) {
	l.zaplog.Error(msg)
	defer l.zaplog.Sync()

}

// Fatal implement logger interface
func (l *ZapAdapter) Fatal(msg string) {
	l.zaplog.Fatal(msg)
	defer l.zaplog.Sync()

}

func (l *ZapAdapter) Debugf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Debugf(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Infof(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Infof(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Warnf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Warnf(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Errorf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Errorf(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Fatalf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Fatalf(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Printf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Infof(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Attach(ctx context.Context) context.Context {
	return Attach(ctx, l)
}

func (l *ZapAdapter) Caller(stack int) Logger {
	return l.clone(l.zaplog.WithOptions(zap.AddCallerSkip(stack)))
}

func (l *ZapAdapter) clone(zaplog *zap.Logger) *ZapAdapter {
	return &ZapAdapter{
		zaplog: zaplog,
	}
}

func getField(key string, val interface{}) zap.Field {
	switch val.(type) {
	case int:
		return zap.Int(key, val.(int))
	case int64:
		return zap.Int64(key, val.(int64))
	case int32:
		return zap.Int32(key, val.(int32))
	case string:
		return zap.String(key, val.(string))
	case fmt.Stringer:
		return zap.Stringer(key, val.(fmt.Stringer))
	case time.Time:
		return zap.Time(key, val.(time.Time))
	case bool:
		return zap.Bool(key, val.(bool))
	case float32:
		return zap.Float32(key, val.(float32))
	case float64:
		return zap.Float64(key, val.(float64))
	case []byte:
		return zap.Binary(key, val.([]byte))
	}
	return zap.Skip()
}
