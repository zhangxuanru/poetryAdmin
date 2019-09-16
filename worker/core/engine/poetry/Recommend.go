/*
@Time : 2019/9/16 16:06
@Author : zxr
@File : recommend
@Software: GoLand
*/
package poetry

import "poetryAdmin/worker/core/grasp/poetry"

//诗词首页推荐信息抓取
type Recommend struct {
}

func NewRecommend() *Recommend {
	return &Recommend{}
}

func (r *Recommend) Run() {
	poetry.NewRecommend().StartGrasp()
}
