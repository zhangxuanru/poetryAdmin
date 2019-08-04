package redis

import (
	"errors"
	redigo "github.com/gomodule/redigo/redis"
)

func RedisPool() *redigo.Pool {
	if G_RedisPool == nil {
		InitRedis("")
	}
	return G_RedisPool
}

//从连接池中取一个连接
func GetConn() (conn redigo.Conn, err error) {
	if G_RedisPool == nil {
		return nil, errors.New("G_RedisPool is nil")
	}
	conn = G_RedisPool.Get()
	if conn == nil {
		return nil, errors.New("conn is nill")
	}
	return conn, nil
}

//发布频道
func Publish(pubTitle string, data string) (reply interface{}, err error) {
	var conn redigo.Conn
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if conn, err = GetConn(); err != nil {
		return nil, err
	}
	if _, err = conn.Do("Publish", pubTitle, data); err != nil {
		return nil, err
	}
	if err = G_RedisPool.Get().Flush(); err != nil {
		return nil, err
	}
	return G_RedisPool.Get().Receive()
}

//获取一个key
func Get(key interface{}) (data interface{}, err error) {
	var conn redigo.Conn
	if conn, err = GetConn(); err != nil {
		return nil, err
	}
	defer conn.Close()
	data, err = conn.Do("GET", key)
	return
}
