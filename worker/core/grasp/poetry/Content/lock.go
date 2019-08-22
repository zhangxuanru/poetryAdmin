package Content

import (
	"poetryAdmin/worker/app/redis"
	"reflect"
)

type Lock struct {
}

func NewLock() *Lock {
	return &Lock{}
}

func (l *Lock) AddKey(key interface{}) bool {
	if _, err := redis.Set(key, 1); err != nil {
		return false
	}
	return true
}

func (l *Lock) DelKey(key interface{}) {
	redis.Del(key)
}

func (l *Lock) ExistsKey(key interface{}) bool {
	var (
		data interface{}
		err  error
	)
	if data, err = redis.Get(key); err != nil {
		return false
	}
	if reflect.TypeOf(data) == nil {
		return false
	}
	if len(data.([]uint8)) > 0 {
		return true
	}
	return false
}