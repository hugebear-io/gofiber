package logger

import (
	"fmt"
	"os"

	"github.com/hugebear-io/gofiber/fabric"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(string, ...zap.Field)
	Info(string, ...zap.Field)
	Warn(string, ...zap.Field)
	Error(error, ...zap.Field)
	Panic(error, ...zap.Field)
	Fatal(error, ...zap.Field)
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
	BuildFields(...interface{}) []zap.Field
	Close()
}

type logger struct {
	logger *zap.Logger
}

func NewLogger(option *LoggerOption) Logger {
	option = validateLoggerOption(option)
	var encoderConfig zapcore.EncoderConfig
	if option.Mode == fabric.DEVELOPMENT_MODE {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.TimeKey = "datetime"
	encoderConfig.LevelKey = "level"
	encoderConfig.CallerKey = "at"
	encoderConfig.MessageKey = "msg"
	var encoder zapcore.Encoder
	if option.JsonEncoding {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	priority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.Level(option.LogLevel)
	})

	core := zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(os.Stdout)), priority)
	if option.Writer != nil && option.Mode == fabric.PRODUCTION_MODE {
		core = zapcore.NewTee(
			core,
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), *option.Writer, priority),
		)
	}

	instance := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(option.SkipCaller))
	return &logger{logger: instance}
}

func (l logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l logger) Error(err error, fields ...zap.Field) {
	l.logger.Error(err.Error(), fields...)
}

func (l logger) Panic(err error, fields ...zap.Field) {
	l.logger.Panic(err.Error(), fields...)
}

func (l logger) Fatal(err error, fields ...zap.Field) {
	l.logger.Fatal(err.Error(), fields...)
}

func (l logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, v...))
}

func (l logger) Infof(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, v...))
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
}

func (l logger) Panicf(format string, v ...interface{}) {
	l.logger.Panic(fmt.Sprintf(format, v...))
}

func (l logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, v...))
}

func (l logger) BuildFields(args ...interface{}) []zap.Field {
	fields := []zap.Field{}
	isEven := len(args)%2 == 0
	if !isEven {
		return fields
	}

	for i := 0; i < len(args); i += 2 {
		key, _ := args[i].(string)
		var field zapcore.Field
		switch v := args[i+1].(type) {
		case string:
			field = zap.String(key, v)
		case int:
			field = zap.Int(key, v)
		case int16:
			field = zap.Int16(key, v)
		case int32:
			field = zap.Int32(key, v)
		case int64:
			field = zap.Int64(key, v)
		case uint:
			field = zap.Uint(key, v)
		case uint16:
			field = zap.Uint16(key, v)
		case uint32:
			field = zap.Uint32(key, v)
		case uint64:
			field = zap.Uint64(key, v)
		case float32:
			field = zap.Float32(key, v)
		case float64:
			field = zap.Float64(key, v)
		case bool:
			field = zap.Bool(key, v)
		default:
			field = zap.Any(key, v)
		}

		fields = append(fields, field)
	}

	return fields
}

func (l *logger) Close() {
	l.logger.Sync()
}
