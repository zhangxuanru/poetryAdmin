package logic

import (
	"poetryAdmin/master/library/config"
	"poetryAdmin/master/library/redis"
	"time"
)

type PublishMsg struct {
	PubTile  string        `json:"pub_tile"`
	AddDate  int64         `json:"add_date"`
	Status   TaskStatus    `json:"status"`
	TaskMark RedisTaskMark `json:"task_mark"`
}

func NewPublishMsg(PubTile string, TaskMark RedisTaskMark) *PublishMsg {
	return &PublishMsg{
		PubTile:  PubTile,
		AddDate:  time.Now().Unix(),
		Status:   TaskStatusImplemented,
		TaskMark: TaskMark,
	}
}

func (p *PublishMsg) BindJson() (string, error) {
	bytes, e := config.G_Json.Marshal(p)
	if e != nil {
		return "", e
	}
	return string(bytes), nil
}

//发布数据
func (p *PublishMsg) PublishData(channel string) (reply interface{}, err error) {
	var data string
	if data, err = p.BindJson(); err != nil {
		return nil, err
	}
	reply, err = redis.Publish(channel, data)
	return reply, err
}
