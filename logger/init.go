package logger

var LoggerInstance Logger = NewLoggerMock()

func InitLogger(option *LoggerOption) {
	LoggerInstance = NewLogger(option)
}

func CloseLogger() {
	LoggerInstance.Close()
}
