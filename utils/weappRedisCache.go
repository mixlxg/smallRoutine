package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type RedisCache struct {
	ctx context.Context
	rdb *redis.Client
	logger *logrus.Logger
}

func NewCache(host string, port int,password string, logger *logrus.Logger) *RedisCache {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:               host+":"+strconv.Itoa(port),
		Password:           password,
		DB:                 0,
	})
	return &RedisCache{
		ctx: ctx,
		rdb: redisClient,
		logger:logger,
	}
}
func(rc *RedisCache) Set(key string, val interface{}, timeout time.Duration)  {
	err:=rc.rdb.Set(rc.ctx,key,val,timeout).Err()
	if err != nil{
		rc.logger.Errorf("redis set key:%s value:%v失败，错误信息：%s",key,val,err.Error())
	}
}

func(rc *RedisCache) Get(key string)(interface{},bool)  {
	val,err := rc.rdb.Get(rc.ctx,key).Result()
	if err != nil{
		rc.logger.Errorf("获取key:%s值失败，错误信息：%s",key,err.Error())
		return nil,false
	}else if err == redis.Nil {
		rc.logger.Errorf("查询Key:%s的值不存在", key)
		return nil,false
	}else {
		rc.logger.Infof("查询key:%s，value:%v成功",key,val)
		return val,true
	}
}
