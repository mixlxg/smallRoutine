package router

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v3"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"smallRoutine/config"
	"smallRoutine/middleware"
	"smallRoutine/views"
)

func NewRouter(config *config.Config,logger *logrus.Logger,gdb *gorm.DB, store redis.Store,cstore base64Captcha.Store, wsdk *weapp.Client,basePath string) (err error)  {
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
	// 配置静态资源
	router.StaticFS(path.Join(config.Http.BaseContext,"statics"),http.Dir(filepath.Join(basePath,config.Http.Static)))
	// 初始化session 中间件
	router.Use(sessions.Sessions(config.Session.SessionId,store))
	// 获取图形验证码的接口
	router.GET(path.Join(config.Http.BaseContext,"/getCaptcha"),views.GetCaptcha(logger,cstore))
	// 登录接口
	router.GET(path.Join(config.Http.BaseContext,"/login"),views.Login(logger, config, wsdk,gdb,cstore))
	// 修改密码接口
	router.POST(path.Join(config.Http.BaseContext,"/app/modifyPwd"),views.ModifyPwd(logger,gdb))

	gr := router.Group(path.Join(config.Http.BaseContext,"/app"),middleware.SmallRoutineSessions(logger))
	{
		// 用户退出接口
		gr.POST("/logout", views.Logout(logger))
		gr.GET("/getUserMess",views.GetUserMess(logger,gdb))
		gr.GET("/getSelectActivityList",views.GetSelectActivityList(logger,gdb))
		gr.GET("/getPageActivityMess",views.GetPageActivityMess(logger,gdb))
		gr.POST("/upload",views.UpLoad(logger))
		//gr.POST("/createOrder",views.CreateOrder(logger,gdb))

	}
	// 管理后台接口
	router.GET(path.Join(config.Http.BaseContext,"/admin/login"),views.AdminLogin(logger,gdb,cstore,config))
	agr := router.Group(path.Join(config.Http.BaseContext,"/admin"),middleware.AdminSessions())
	{
		agr.POST("/logout",views.AdminLogout(logger))
		agr.GET("/getCompany",views.GetCompany(logger,gdb))
		agr.POST("/queryUser", views.QueryUser(logger,gdb))
		agr.POST("/pageQueryUsers",views.PageQueryUser(logger,gdb))
		agr.POST("/createUser",views.CreateUser(logger,gdb))
		agr.GET("/getCurrentUserMess",views.GetCurrentUserMess(logger,gdb))
		agr.POST("/updateUser",views.UpdateUser(logger,gdb))
		agr.GET("/delUser",views.DelUser(logger,gdb))
		agr.POST("/createActivity",views.CreateActivity(logger,gdb))
		agr.POST("/MdApprover",views.MdApprover(logger,gdb))
		agr.GET("/queryActivity",views.QueryActivityMess(logger,gdb))
		agr.GET("/delActivity",views.DelActivity(logger,gdb))
		agr.POST("/updateActivity",views.UpdateActivity(logger,gdb))
		agr.POST("/createGroup",views.CreateGroup(logger,gdb))
		agr.POST("/addUsersToGroup",views.AddUsersToGroup(logger,gdb))
		agr.POST("/delUserFromGroup",views.DelUserFromGroup(logger,gdb))
		agr.POST("/delGroup",views.DelGroup(logger,gdb))
		agr.POST("/modifyGroup",views.GroupModify(logger,gdb))
		agr.POST("/setGroupLeader",views.SetGroupLeader(logger,gdb))
		agr.POST("/queryActivityGroupsUsers",views.QueryActivity(logger,gdb))

	}
	//绑定地址端口启动服务
	err = router.Run(fmt.Sprintf("%s:%d",config.Http.Host,config.Http.Port))
	return
}
