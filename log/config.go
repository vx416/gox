package log

import (
	"errors"
	"strings"

	"github.com/vx416/gox/log/encoder"
	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
)

var (
	// ErrUnknownLoggerType represent unknown logger type  error
	ErrUnknownLoggerType = errors.New("logger type unknown")
)

// LoggerType define logger engine type
type LoggerType string

const (
	// Zap represent zap logger engine
	Zap LoggerType = "zap"
)

// Env define logger environment
type Env string

// IsDev make env to lower and check if string is dev
func (env Env) IsDev() bool {
	envStr := strings.ToLower(string(env))

	return envStr == "dev" || envStr == "development"
}

// Level define logger level
type Level string

// ZapLevel convert string to zapcore.Level
func (l Level) ZapLevel() zapcore.Level {
	var (
		zapLevel zapcore.Level
		level    = strings.ToLower(string(l))
	)
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	case "panic":
		zapLevel = zapcore.PanicLevel
	}
	return zapLevel
}

// Config represetn logger config
type Config struct {
	AppName        string     `yaml:"app_name"`
	Env            Env        `yaml:"env"`
	Level          Level      `yaml:"level"`
	OutputPaths    []string   `yaml:"output_paths"`
	ErrOutputPaths []string   `yaml:"err_output_paths"`
	Type           LoggerType `yaml:"type"`
}

// Build build logger
func (config Config) Build() (Logger, error) {
	var (
		logger Logger
		err    error
	)

	switch config.Type {
	case Zap:
		logger, err = config.buildZap()
	default:
		err = ErrUnknownLoggerType
	}
	if err != nil {
		return nil, err
	}

	SetGlobal(logger)
	return logger, nil
}

func (config Config) buildZap() (Logger, error) {
	var (
		zapConfig zap.Config
	)
	if config.Env.IsDev() {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(config.Level.ZapLevel())
		zapConfig.EncoderConfig.EncodeLevel = encoder.ColorfulLevelEncoder
		zapConfig.EncoderConfig.EncodeCaller = encoder.ColorizeCallerEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	if len(config.OutputPaths) > 0 {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, config.OutputPaths...)
	}
	if len(config.ErrOutputPaths) > 0 {
		zapConfig.ErrorOutputPaths = append(zapConfig.OutputPaths, config.ErrOutputPaths...)
	}

	field := zap.String("app_name", config.AppName)

	zaplog, err := zapConfig.Build(zap.Fields(field), zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &ZapAdapter{zaplog: zaplog, env: config.Env}, nil
}
