package redis

import (
	"errors"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"time"
)

var G_RedisPool *redigo.Pool

func InitRedis(addr string) (err error) {
	if addr == "" {
		return errors.New("addr is nil")
	}
	logrus.Debug("redis addr:", addr)
	var dial redigo.Conn
	G_RedisPool = &redigo.Pool{
		Dial: func() (conn redigo.Conn, e error) {
			dial, err = redigo.Dial("tcp", addr)
			if err != nil {
				logrus.Debug("redis dial error:", err)
				return nil, err
			}
			return dial, err
		},
		MaxIdle:     10,                               //最初的连接数量，池子里的最大空闲连接
		MaxActive:   0,                                //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: time.Duration(time.Second * 180), //超过这个duration的空闲连接，会被关闭
		Wait:        true,
	}
	return
}
