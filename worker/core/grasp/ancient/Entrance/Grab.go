/*
@Time : 2019/8/30 17:34
@Author : zxr
@File : Grab
@Software: GoLand
*/
package Entrance

import "poetryAdmin/worker/core/grasp/ancient/Action"

type Grab struct {
}

func NewGrab() *Grab {
	return &Grab{}
}

//开始执行古文抓取
func (g *Grab) Exec() {
	//先从首页抓
	Action.NewIndex().StartGrab()
}
