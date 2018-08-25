package redisdb

import (
	"github.com/garyburd/redigo/redis"
	"DBServer/conf"
	"time"
	"github.com/cihub/seelog"
)

var RedisPool *redis.Pool

func Init(){
	RedisPool = &redis.Pool{
		MaxIdle:conf.ServerData.RedisMaxIdleNum,
		MaxActive:conf.ServerData.RedisMaxActiveNum,
		IdleTimeout: time.Duration(conf.ServerData.RedisIdleTimeOut) * time.Second,
		Dial: func() (redis.Conn, error) {
			do := redis.DialPassword(conf.ServerData.RedisPass)
			return redis.Dial("tcp", conf.ServerData.RedisAddr, do)
		},
	}
	seelog.Debugf("connet redis success")
}
