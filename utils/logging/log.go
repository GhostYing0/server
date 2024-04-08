package logging

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// L 日志实例
var L *zap.SugaredLogger

// Setup initialize the log instance
func Setup() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	L = logger.Sugar()
}

func New(path string) *zap.Logger {
	writeSyncer := getLogWriter(path)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	return logger
}

// 日志编码方式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.00")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 日志打印指定路径及配置
func getLogWriter(path ...string) zapcore.WriteSyncer {
	logPath := "./logs/test.log"
	if len(path) > 0 {
		logPath = path[0]
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath, // 日志文件路径
		MaxSize:    1,       // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 10,      // 保留旧文件的最大个数
		MaxAge:     30,      // 保留旧文件的最大天数
		Compress:   false,   // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}
