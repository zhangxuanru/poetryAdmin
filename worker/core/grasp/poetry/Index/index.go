package Index

//抓取首页
type Index struct {
	content chan string
}

func NewIndex() *Index {
	return &Index{}
}

func (i *Index) GetAllData() {

}

//首页-诗文分类
func (i *Index) GetPoetryCategory() {

}

//首页-名句分类
func (i *Index) GetPoetryFamousCategory() {

}

//首页-作者
func (i *Index) GetPoetryAuthor() {

}

//获取首页html内容
func (i *Index) GetIndexHtml() {

}
