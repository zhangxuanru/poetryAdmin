package data

//保存诗词正文和赏析注释结果...
type contentStore struct {
}

func NewContentStore() *contentStore {
	return new(contentStore)
}

func (c *contentStore) LoadPoetryContentData(data interface{}, params interface{}) {

}
