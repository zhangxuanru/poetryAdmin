package define

//
type HomeFormat struct {
	Identifier DataFormat
	Data       interface{}
}

//首页保存的数据格式
type TextHrefFormat struct {
	Text         string
	Href         string
	ShowPosition ShowPosition
}

//分类页-作者与诗文对应关系【体裁下对应着多个诗文】
type PoetryAuthorList struct {
	AuthorName      string //作者姓名
	PoetryTitle     string //诗词标题
	PoetrySourceUrl string //诗词链接
	GenreTitle      string //体裁名称
	Category        *TextHrefFormat
}
