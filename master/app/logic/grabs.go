package logic

import (
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

//发布数据
func (p *PublishMsg) PublishData(channel string, data string) (reply interface{}, err error) {
	reply, err = redis.Publish(channel, data)
	return reply, err
}
