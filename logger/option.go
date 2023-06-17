package logger

import (
	"github.com/hugebear-io/gofiber/fabric"
	"go.uber.org/zap/zapcore"
)

type LoggerOption struct {
	Mode         string
	LogLevel     LogLevel
	SkipCaller   int
	JsonEncoding bool
	Writer       *zapcore.WriteSyncer
}

var LoggerDefaultOption LoggerOption = LoggerOption{
	Mode:         fabric.DEVELOPMENT_MODE,
	LogLevel:     LOG_LEVEL_DEBUG,
	JsonEncoding: false,
	SkipCaller:   1,
	Writer:       nil,
}

func validateLoggerOption(opt *LoggerOption) *LoggerOption {
	if opt == nil {
		return &LoggerDefaultOption
	}

	return opt
}
