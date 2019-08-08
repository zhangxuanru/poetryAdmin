package poetry

import (
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Index"
	"reflect"
)

//诗词全站抓取
type RunAll struct {
}

func NewRunAll() *RunAll {
	return &RunAll{}
}

//执行全站抓取
func (r *RunAll) Run() {
	//先获取锁  临时注释
	//if _, err := redis.SetNx(r.GetLockKey(), "1", "3600"); err != nil {
	//	go data.G_GraspResult.PushErrorAndClose(err)
	//	return
	//}
	r.Execution()
	return
}

//获取锁的键
func (r *RunAll) GetLockKey() (key string) {
	//key := (string)(define.RedisIsTaskAllRun)
	key = reflect.ValueOf(define.RedisIsTaskAllRun).String()
	return key
}

//执行抓取
func (r *RunAll) Execution() {
	Index.NewIndex().GetAllData()
	defer data.G_GraspResult.PushCloseMark(true)
}
