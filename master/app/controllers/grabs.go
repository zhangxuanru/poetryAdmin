package controllers

import (
	"net/http"
	"poetryAdmin/master/app/logic"
	"poetryAdmin/master/library/config"
	"poetryAdmin/master/library/redis"
)

//一键抓取列表
func Grabs(w http.ResponseWriter, r *http.Request) {
	isRun, _ := redis.Get(logic.RedisIsTaskAllRun)
	if isRun == nil {
		isRun = ""
	}
	ret := make(map[string]interface{})
	ret["is_run"] = isRun.(string)
	base.DisplayHtmlLayOut(w, "grab-list.html", ret, nil)
}

//执行抓取, 写入redis
func GrabsImpl(w http.ResponseWriter, r *http.Request) {
	var (
		jsonData []byte
		err      error
	)
	publicMsg := logic.NewPublishMsg(logic.GrabTaskTitleAll, logic.GrabPoetryAll)
	if jsonData, err = config.G_Json.Marshal(publicMsg); err != nil {
		goto OutPutERR
	}
	if _, err = publicMsg.PublishData(config.G_Conf.PubChannelTitle, string(jsonData)); err != nil {
		goto OutPutERR
	}
	base.OutPutRespJson(w, nil, logic.GrabTaskAdd, logic.RespSuccess)
	return
OutPutERR:
	if err != nil {
		base.OutPutRespJson(w, nil, err.Error(), logic.RespFail)
	} else {
		base.OutPutRespJson(w, nil, logic.RespFailMsg, logic.RespFail)
	}
	return
}

//抓取结果列表页
func GrabsList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("抓取结果列表页"))
}
