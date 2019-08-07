package test

import (
	. "github.com/smartystreets/goconvey/convey"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/parse"
	"testing"
	"time"
)

func TestGrabsIndex(t *testing.T) {
	var SubReceiveMsg parse.SubscribeMsg
	Convey("测试抓取所有", t, func() {
		SubReceiveMsg = parse.SubscribeMsg{
			PubTile:  "抓取所有",
			AddDate:  time.Now().Unix(),
			Status:   define.TaskStatusImplemented,
			TaskMark: define.GrabPoetryAll,
		}
		parse.NewDispatch(SubReceiveMsg).Execution()

		time.Sleep(10 * time.Second)
	})

}
