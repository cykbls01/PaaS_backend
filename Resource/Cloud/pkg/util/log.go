package util

import (
	"Cloud/pkg/consts"
	log "github.com/wzyonggege/logger"
)

//var logger *zap.Logger

//func init(){
//	hook := lumberjack.Logger{
//		Filename: 		consts.PATH + "logs/paas.log",
//		MaxSize: 		128,
//		MaxBackups:		30,
//		MaxAge:			7,
//		Compress:		false,
//	}
//
//	encoderConfig := zapcore.EncoderConfig{
//		TimeKey:		"time",
//		LevelKey:		"level",
//		NameKey:		"logger",
//		//CallerKey:		"linenum",
//		MessageKey:		"msg",
//		StacktraceKey:	"stacktrace",
//		LineEnding:		zapcore.DefaultLineEnding,
//		EncodeLevel:	zapcore.LowercaseLevelEncoder,
//		EncodeTime:		zapcore.ISO8601TimeEncoder,
//		EncodeDuration:	zapcore.SecondsDurationEncoder,
//		EncodeCaller:	zapcore.FullCallerEncoder,
//		EncodeName:		zapcore.FullNameEncoder,
//	}
//
//	atomicLevel := zap.NewAtomicLevel()
//	atomicLevel.SetLevel(zap.InfoLevel)
//
//	core := zapcore.NewCore(
//		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
//		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
//		atomicLevel,                                                                     // 日志级别
//	)
//
//	// 开启开发模式，堆栈跟踪
//	caller := zap.AddCaller()
//	// 开启文件及行号
//	development := zap.Development()
//	// 构造日志
//	logger = zap.New(core, caller, development)
//}

func init(){
	c := log.New()
	c.SetDivision("time")
	c.SetTimeUnit(log.Day)
	c.SetEncoding("json")

	c.SetInfoFile(consts.PATH + "logs/paas.log")
	c.SetErrorFile(consts.PATH + "logs/paas_err.log")

	c.InitLogger()
}
//
//func GetLogger() *zap.Logger{
//	return logger
//}
