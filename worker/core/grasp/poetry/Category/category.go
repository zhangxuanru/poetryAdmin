package Category

import "poetryAdmin/worker/core/define"

//诗文分类模块 抓取诗文分类
type Category struct {
}

func NewCategory() *Category {
	return new(Category)
}

//通过首页抓取到的诗文分类传到这里，这里循环数据去发送请求
func (c *Category) GraspByIndexData(data *define.HomeFormat) {

}
