/*
@Time : 2019/9/2 14:37
@Author : zxr
@File : ancient_test
@Software: GoLand
*/
package test

import (
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/grasp/ancient/Entrance"
	"testing"
	"time"
)

//测试古籍
func TestAncient(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	Entrance.NewGrab().Exec()
	time.Sleep(60 * time.Second)
}
