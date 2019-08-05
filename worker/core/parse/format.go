package parse

import (
	"poetryAdmin/worker/core/define"
)

type SubscribeMsg struct {
	PubTile  string               `json:"pub_tile"`
	AddDate  int64                `json:"add_date"`
	Status   define.TaskStatus    `json:"status"`
	TaskMark define.RedisTaskMark `json:"task_mark"`
}

func NewSubscribeMsg() *SubscribeMsg {
	return new(SubscribeMsg)
}
