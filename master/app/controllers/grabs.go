package controllers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"poetryAdmin/master/app/logic"
	"poetryAdmin/master/library/config"
	"poetryAdmin/master/library/redis"
)

//一键抓取列表
func Grabs(w http.ResponseWriter, r *http.Request) {
	data, _ := redis.Get(logic.RedisIsTaskAllRun)
	ret := make(map[string]interface{})
	ret["redisData"] = data.(string)
	base.DisplayHtmlLayOut(w, "grab-list.html", ret, nil)
}

//执行抓取, 写入redis
func GrabsImpl(w http.ResponseWriter, r *http.Request) {
	publicMsg := logic.NewPublishMsg(logic.GrabTaskTitleAll, logic.GrabPoetryAll)
	reply, err := publicMsg.PublishData(config.G_Conf.PubChannelTitle)
	if err != nil {
		base.OutPutRespJson(w, nil, err.Error(), logic.RespFail)
		return
	}
	logrus.Info("reply:", reply)
	base.OutPutRespJson(w, nil, logic.GrabTaskAdd, logic.RespSuccess)
	return
}
