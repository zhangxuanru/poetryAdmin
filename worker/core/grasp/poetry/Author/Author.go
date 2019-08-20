package Author

import "poetryAdmin/worker/core/define"

//作者模块  抓取作者
type Author struct {
}

func NewAuthor() *Author {
	return new(Author)
}

//通过首页抓取到的作者列表传到这里，这里循环数据去发送请求
func (a *Author) GraspByIndexData(data *define.HomeFormat) {

}

//获取作者详情信息
func (a *Author) GetAuthorDetail(authorUrl string) {

}

//获取作者诗词列表数据，并保存诗词列表
func (a *Author) GetAuthorPoetryList(authorUrl string) {

}
