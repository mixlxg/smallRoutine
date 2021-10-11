package utils

import (
	wlog "github.com/medivhzhan/weapp/v3/logger"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"smallRoutine/config"
)

// 定义一个logger 初始化的方法
func NewLogger(config *config.LogrusConfig) *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		DisableColors:			   true,
		DisableQuote:              true,
		FullTimestamp:			   true,
	})
	//拼接日志文件名称和路径
	logFileName := filepath.Join(config.LogPath,config.AppLogName)
	// 初始化一个lumberjack logger
	logFileObj := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    config.FileMaxSize,
		LocalTime:  true,
	}
	/* logFileObj, err:= os.OpenFile(logFileName,os.O_CREATE|os.O_RDWR,0644)
	if err !=nil{
		panic(fmt.Sprintf("打开日志文件失败，错误信息：%v",err))
	}*/
	// 定义一个io.writer 数组用于存放多个io.writer接口
	var writers []io.Writer = []io.Writer{logFileObj}
	// 判断日志是否输出到控制台
	if config.Console{
		writers = append(writers,os.Stdout)
	}
	logger.SetReportCaller(true)
	logger.SetOutput(io.MultiWriter(writers...))
	// 绑定配置日志的级别
	switch config.Level {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	case "fatal":
		logger.SetLevel(log.FatalLevel)
	case "panic":
		logger.SetLevel(log.PanicLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}
	return logger
}

// 初始化小程序sdk 日志handle
func NewWexinLogegr(config *config.LogrusConfig) wlog.Logger  {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		DisableColors:			   true,
		DisableQuote:              true,
		FullTimestamp:			   true,
	})
	//拼接日志文件名称和路径
	logFileName := filepath.Join(config.LogPath,config.WeixinLog)
	// 初始化一个lumberjack logger
	logFileObj := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    config.FileMaxSize,
		LocalTime:  true,
	}
	// 定义一个io.writer 数组用于存放多个io.writer接口
	var writers []io.Writer = []io.Writer{logFileObj}
	// 判断日志是否输出到控制台
	if config.Console{
		writers = append(writers,os.Stdout)
	}
	logger.SetReportCaller(true)
	logger.SetOutput(io.MultiWriter(writers...))
	// 绑定配置日志的级别
	switch config.Level {
	case "debug":
		return wlog.NewLogger(logger,wlog.Info,false)
	case "info":
		return wlog.NewLogger(logger,wlog.Info,false)
	case "warn":
		return wlog.NewLogger(logger,wlog.Warn,false)
	case "error":
		return wlog.NewLogger(logger,wlog.Error,false)
	case "fatal":
		return wlog.NewLogger(logger,wlog.Error,false)
	case "panic":
		return wlog.NewLogger(logger,wlog.Error,false)
	default:
		return wlog.NewLogger(logger,wlog.Info,false)
	}
}



