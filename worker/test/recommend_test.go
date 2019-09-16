/*
@Time : 2019/9/16 16:15
@Author : zxr
@File : recommend_test
@Software: GoLand
*/
package test

import (
	"poetryAdmin/worker/core/grasp/poetry"
	"testing"
)

//抓取推荐信息
func TestRecommend(t *testing.T) {
	poetry.NewRecommend().StartGrasp()
}
