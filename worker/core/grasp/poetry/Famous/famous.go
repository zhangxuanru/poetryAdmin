package Famous

import "poetryAdmin/worker/core/define"

//名句模块  抓取名句
type Famous struct {
}

func NewFamous() *Famous {
	return new(Famous)
}

//通过首页抓取到的名句分类传到这里，这里循环数据去发送请求
func (f *Famous) GraspByIndexData(data *define.HomeFormat) {

}
