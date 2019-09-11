package poetry

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/redis"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/ancient/Entrance"
	Entrance2 "poetryAdmin/worker/core/grasp/famous/Entrance"
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
	defer func() {
		_, _ = redis.Del(r.GetLockKey())
	}()
	//先获取锁  临时注释
	if config.G_Conf.Env != define.TestEnv {
		if _, err := redis.SetNx(r.GetLockKey(), "1", "3600"); err != nil {
			logrus.Infoln("err:", err)
			go data.G_GraspResult.PushErrorAndClose(err)
			return
		}
	}
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
	//抓取古籍
	go Entrance.NewGrab().Exec()
	//抓取名句
	go Entrance2.NewFamous().Run()
	//抓取诗词
	Index.NewIndex().GetAllData()

	logrus.Infoln("结果处理结束......")

	//临时关掉， 还没确定在哪一步关闭获取结果的goroutine
	//defer data.G_GraspResult.PushCloseMark(true)
}
