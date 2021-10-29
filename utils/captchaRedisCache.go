package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type CaptchaRedisCache struct {
	ctx context.Context
	rdb	*redis.Client
	logger *logrus.Logger
	timeout	time.Duration
}

func NewCaptchaRedisCache(host string,port int,password string,timeout time.Duration,logger *logrus.Logger) *CaptchaRedisCache {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:               host+":"+strconv.Itoa(port),
		Password:           password,
		DB:                 0,
	})
	return &CaptchaRedisCache{
		ctx: ctx,
		rdb: redisClient,
		logger:logger,
		timeout:timeout,
	}
}
func(crc *CaptchaRedisCache) Set(id string, value string) error {
	err:=crc.rdb.Set(crc.ctx,id,value,time.Second*crc.timeout).Err()
	if err !=nil{
		crc.logger.Errorf("图形验证码设置redis失败key:%s,value:%s,错误信息：%s",id,value,err.Error())
	}
	return err
}

func(crc *CaptchaRedisCache) Get(id string, clear bool) string  {
	res,err := crc.rdb.Get(crc.ctx,id).Result()
	if err !=nil{
		crc.logger.Errorf("图形验证码获取redis，key:%s值失败，错误信息：%s",id,err.Error())
		return ""
	}else if err == redis.Nil{
			crc.logger.Infof("图形验证码获取redis，key:%s值不存在",id)
			return ""
	}
	if clear{
		err = crc.rdb.Del(crc.ctx,id).Err()
		if err !=nil{
			crc.logger.Warnf("图形验证码删除key:%s失败，错误信息：%s",id,err.Error())
		}
	}
	return res
}

func(crc *CaptchaRedisCache) Verify(id, answer string, clear bool)bool  {
	res,err := crc.rdb.Get(crc.ctx,id).Result()
	if err !=nil{
		crc.logger.Errorf("图形验证码获取redis,key:%s失败，错误信息：%s",id,err.Error())
		return false
	}
	if clear{
		err = crc.rdb.Del(crc.ctx,id).Err()
		if err !=nil{
			crc.logger.Warnf("图形验证码删除key:%s失败，错误信息：%s",id,err.Error())
		}
	}
	if res == answer{
		return true
	}else {
		return false
	}
}
