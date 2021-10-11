package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gin-contrib/sessions/redis"
	"github.com/medivhzhan/weapp/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"smallRoutine/config"
	"smallRoutine/model"
	"smallRoutine/utils"
	"strconv"
	"strings"
)
var basePath string
var conf *config.Config
var logger *logrus.Logger
var gdb *gorm.DB
var err error
var store redis.Store
var wSdk *weapp.Client
func init()  {
	// 获取base path
	basePath,err = os.Getwd()
	if err !=nil{
		panic(fmt.Sprintf("获取项目的BasePath路径失败，错误信息:%s",err.Error()))
	}
	// 如果是编译后执行，二进制文件应该存放在bin目录下这时候获取到的os.getwd就是bin的路径了
	//所以需要处理一下
	if strings.Contains(basePath,"bin"){
		basePath = filepath.Dir(basePath)
	}
	// 拼接applicaton.yaml 绝对路径
	configPath := filepath.Join(basePath,"application.yml")
	// 解析获取application.yaml 内容
	conf,err = config.NewConfig(configPath)
	if err != nil{
		panic(err)
	}
	// 初始化日志对象
	logger = utils.NewLogger(conf.Logrus)
	logger.Info("初始化日志成功")
	// 初始化数据库链接
	logger.Info("开始初始化数据库链接")
	gdb,err = utils.NewGorm(conf)
	if err != nil{
		logger.Fatalf("初始化%s数据库报错，错误信息：%v",conf.Mysql.MysqlUrl,err)
	}
	logger.Info("初始化数据库成功")
	// 获取命令行是否有传初始化数据库表的参数
	isInitTable := flag.Bool("initDb",false,"初始化，并创建系统表")
	flag.Parse()
	if *isInitTable {
		model.Init(gdb)
	}
	store,err = redis.NewStore(conf.Redis.MaxIdle,"tcp",conf.Redis.Host+":"+strconv.Itoa(conf.Redis.Port),conf.Redis.Password,[]byte(conf.Session.Secret))
	if err !=nil{
		logger.Errorf("初始化redis session 存储失败，报错信息：%s", err.Error())
		panic(err)
	}
	// 初始化一个微信小程序的sdk
	//1.自定义一个http客户端
	cli := &http.Client{
		Transport:     &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
		},
		Timeout:       conf.WeixinConfig.HTimeOut,
	}
	// 2. 日志直接使用初始化好的logger
	wlog := utils.NewWexinLogegr(conf.Logrus)
	// 自定义缓存
	cache := utils.NewCache(conf.Redis.Host,conf.Redis.Port,conf.Redis.Password,logger)
	// 自定义token获取方式
	wSdk = weapp.NewClient(
		conf.WeixinConfig.AppId,
		conf.WeixinConfig.Secret,
		weapp.WithHttpClient(cli),
		weapp.WithLogger(wlog),
		weapp.WithCache(cache),
		weapp.WithAccessTokenSetter(utils.GetToken(conf,logger)),
		)
}
