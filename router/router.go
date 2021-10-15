package router

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"io"
	"os"
	"path"
	"path/filepath"
	"smallRoutine/config"
	"smallRoutine/middleware"
	"smallRoutine/views"
)

func NewRouter(config *config.Config,logger *logrus.Logger,gdb *gorm.DB,basepath string, store redis.Store, wsdk *weapp.Client) (err error)  {
	//初始化gin配置
	// 关闭gin日志着色配置
	gin.DisableConsoleColor()
	// 绑定默认日志wirter
	var writers []io.Writer
	// 判断是否需要在控制台输出日志，并初始化
	if config.Logrus.Console{
		writers = append(writers, os.Stdout)
	}
	//初始化一个lumberjack的writer 实现gin日志的按照配置切割
	writers = append(writers,&lumberjack.Logger{
		Filename:   filepath.Join(config.Logrus.LogPath,config.Logrus.GinLogName),
		MaxSize:    config.Logrus.FileMaxSize,
		LocalTime:  true,
	})
	//绑定gin默认writer
	gin.DefaultErrorWriter = io.MultiWriter(writers...)
	gin.DefaultWriter = io.MultiWriter(writers...)
	// 初始化一个默认gin 引擎
	router := gin.New()
	// 添加日志中间件,配置日志格式
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] %s %s %s %d %s %s %s %d\n",
			param.ClientIP,
			param.TimeStamp.String(),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
			param.BodySize)
	}))
	// 添加接口奔溃恢复中间件
	router.Use(gin.Recovery())
	// 初始化session 中间件
	router.Use(sessions.Sessions(config.Session.SessionId,store))
	// 登录接口
	router.GET(path.Join(config.Http.BaseContext,"/login"),views.Login(logger, config, wsdk,gdb))
	gr := router.Group(path.Join(config.Http.BaseContext,"/app"),middleware.SmallRoutineSessions(logger))
	{
		// 用户退出接口
		gr.POST("/logout", views.Logout(logger))
		// 修改密码接口
		gr.POST("/modifyPwd",views.ModifyPwd(logger,gdb))


	}
	//绑定地址端口启动服务
	err = router.Run(fmt.Sprintf("%s:%d",config.Http.Host,config.Http.Port))
	return
}
