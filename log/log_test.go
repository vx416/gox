package log

import (
	"context"
	stderr "errors"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var config = Config{
	AppName: "test_app",
	Env:     "dev",
	Level:   "debug",
}

func TestZapLog(t *testing.T) {
	config.Type = "zap"

	logger, err := config.Build()
	assert.NoError(t, err)
	logger.Debug("debug test")
	logger2 := logger.Fields(map[string]interface{}{"test_field": 123})
	ctx := context.Background()
	ctx = Attach(ctx, logger2)
	logger.Debug("without field")

	logger2 = Ctx(ctx)
	assert.NoError(t, err)
	logger2.Debug("with field")

	getError(logger)
	logger.Err(getNormalErr()).Error("normal error!!!")
}

func getError(logger Logger) error {
	logger.Caller(1).Error("error!!!")
	return errors.WithStack(errors.New("test error"))
}

func getNormalErr() error {
	return stderr.New("test normal error")
}
