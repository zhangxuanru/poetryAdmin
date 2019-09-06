/*
@Time : 2019/9/6 19:07
@Author : zxr
@File : famous_test
@Software: GoLand
*/
package test

import (
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/grasp/famous/Action"
	"testing"
	"time"
)

//抓取名句首页
func TestFamousIndex(t *testing.T) {
	go data.NewGraspResult().PrintMsg()

	Action.NewIndex().Start()

	time.Sleep(5 * time.Second)
}
