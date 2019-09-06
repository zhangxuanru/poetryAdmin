/*
@Time : 2019/9/6 18:39
@Author : zxr
@File : famous
@Software: GoLand
*/
package Entrance

import "poetryAdmin/worker/core/grasp/famous/Action"

//抓取名句
type famous struct {
}

func NewFamous() *famous {
	return &famous{}
}

func (f *famous) Run() {
	Action.NewIndex().Start()
}
