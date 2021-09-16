package utils

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"fmt"
	"smallRoutine/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

//初始化一个gorm引擎

func NewGorm(config *config.Config) (*gorm.DB,error) {
	// 拼接一个mysql数据库链接url
	dns := fmt.Sprintf("%s:%s@%s",
		config.Mysql.Username,
		config.Mysql.Password,
		config.Mysql.MysqlUrl)
	// 初始化logger
	//拼接日志文件名称和路径
	logFileName := filepath.Join(config.Logrus.LogPath,config.Logrus.SqlLogName)
	// 初始化一个lumberjack logger
	logFileObj := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    config.Logrus.FileMaxSize,
		LocalTime:  true,
	}
	/* logFileObj, err:= os.OpenFile(logFileName,os.O_CREATE|os.O_RDWR,0644)
	if err !=nil{
		panic(fmt.Sprintf("打开日志文件失败，错误信息：%v",err))
	}*/
	// 定义一个io.writer 数组用于存放多个io.writer接口
	var writers []io.Writer = []io.Writer{logFileObj}
	// 判断日志是否输出到控制台
	if config.Logrus.Console{
		writers = append(writers,os.Stdout)
	}
	lger := logger.New(log.New(io.MultiWriter(writers...),"",log.Llongfile|log.LstdFlags),
		logger.Config{
		SlowThreshold: 3*time.Second,
		Colorful:      false,
		LogLevel:      logger.Info,
	}) 
	gormDb,err:=gorm.Open(
		mysql.New(mysql.Config{
			DSN: dns,
		}),
		&gorm.Config{
			Logger:lger,
			CreateBatchSize: 5000,
		})
	if err ==nil {
		db,err := gormDb.DB()
		if err !=nil{
			return nil,err
		}else{
			db.SetConnMaxLifetime(time.Duration(config.Mysql.MaxLifeTime))
			db.SetConnMaxIdleTime(time.Duration(config.Mysql.MaxIdleTime))
			db.SetMaxIdleConns(int(config.Mysql.MaxIdleCount))
			db.SetMaxOpenConns(int(config.Mysql.MaxOpen))
		}
	}
	return gormDb,err
}
