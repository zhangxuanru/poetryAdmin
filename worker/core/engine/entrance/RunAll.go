package entrance

import (
	"poetryAdmin/worker/app/redis"
	"poetryAdmin/worker/core/define"
	"reflect"
)

//全站抓取
type RunAll struct {
}

func NewRunAll() *RunAll {
	return &RunAll{}
}

//执行全站抓取
func (r *RunAll) Run() (err error) {
	//key := (string)(define.RedisIsTaskAllRun)
	key := reflect.ValueOf(define.RedisIsTaskAllRun).String()
	if _, err = redis.SetNx(key, "1", "3600"); err != nil {
		return err
	}

	return
}
