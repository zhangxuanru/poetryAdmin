package redis

import (
	"errors"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"reflect"
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
	if err = conn.Flush(); err != nil {
		return nil, err
	}
	return
}

//订阅频道
func SubScribe(chanTitle string, receive chan []byte) (err error) {
	var conn redigo.Conn
	if conn, err = GetConn(); err != nil {
		logrus.Debug("redis err:", err)
		return err
	}
	defer conn.Close()
	subConn := redigo.PubSubConn{Conn: conn}
	if err = subConn.Subscribe(chanTitle); err != nil {
		return
	}
	for {
		switch v := subConn.Receive().(type) {
		case redigo.Message:
			//logrus.Infof("%s: message: %s\n", v.Channel, v.Data)
			receive <- v.Data
		case redigo.Subscription:
			//logrus.Infof("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return v
		}
	}
	return
}

//获取一个key
func Get(key interface{}) (data interface{}, err error) {
	var conn redigo.Conn
	if conn, err = GetConn(); err != nil {
		logrus.Debug("redis err:", err)
		return nil, err
	}
	defer conn.Close()
	data, err = conn.Do("GET", key)
	if err != nil {
		logrus.Debug("redis err:", err)
	}
	return
}

//设置一个key
func Set(args ...interface{}) (reply interface{}, err error) {
	var conn redigo.Conn
	if conn, err = GetConn(); err != nil {
		logrus.Debug("redis err:", err)
		return
	}
	defer conn.Close()
	reply, err = conn.Do("SET", args...)
	//c.Do("SET", "mykey", "superWang", "EX", "5","NX")
	return
}

//删除一个key
func Del(args ...interface{}) (reply interface{}, err error) {
	var conn redigo.Conn
	if conn, err = GetConn(); err != nil {
		logrus.Debug("redis err:", err)
		return
	}
	defer conn.Close()
	reply, err = conn.Do("DEL", args...)
	return
}

//设置一个key，只有不存在时才设置成功
func SetNx(key string, val string, expire string) (lock bool, err error) {
	var (
		reply  interface{}
		locRet string
	)
	if reply, err = Set(key, val, "EX", expire, "NX"); err != nil {
		return false, err
	}
	if reply == nil {
		return false, errors.New("设置失败")
	}
	locRet = reflect.ValueOf(reply).String()
	if locRet != "OK" || locRet != "ok" {
		return false, errors.New("设置失败")
	}
	return true, nil
}
