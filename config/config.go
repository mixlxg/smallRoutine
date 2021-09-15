package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// mysql配置设计到的字段结构体
type MysqlConfig struct {
	MysqlUrl	string				`yaml:"url"`
	Username 	string				`yaml:"username"`
	Password	string				`yaml:"password"`
	MaxLifeTime	int64				`yaml:"maxlifetime"`
	MaxIdleTime	int64				`yaml:"maxidletime"`
	MaxIdleCount int64				`yaml:"maxidlecount"`
	MaxOpen		int64				`yaml:"maxopen"`
}
// http服务启动设计的配置属性结构体
type HttpConfig struct {
	Host 		string				`yaml:"host"`
	Port		int64				`yaml:"port"`
	BaseContext	string				`yaml:"context"`
}
// zap 日志参数
type LogrusConfig struct {
	LogPath	string					`yaml:"logpath"`
	AppLogName string				`yaml:"applogname"`
	GinLogName string				`yaml:"ginlogname"`
	SqlLogName string				`yaml:"sqllogname"`
	Level	string					`yaml:"level"`
	Console bool					`yaml:"console"`
	FileMaxSize int					`yaml:"filemaxsize"`
}

type Config struct {
	Mysql *MysqlConfig					`yaml:"mysql"`
	Http *HttpConfig					`yaml:"http"`
	Logrus *LogrusConfig				`yaml:"logrus"`
}

func NewConfig(configPath string) (*Config, error)  {
	// 定义一个初始config
	var  config Config
	configContent,err := ioutil.ReadFile(configPath)
	if err != nil {
		return &config, errors.New(fmt.Sprintf("读取配置文件%s失败，错误信息：%s", configPath,err.Error()))
	}
	if err=yaml.Unmarshal(configContent, &config); err !=nil{
		return &config, errors.New(fmt.Sprintf("解析yaml文件失败错误信息：%s", err.Error()))
	}
	return &config,nil
}