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

//诗词详情
type PoetryContent struct {
	Title        string
	Content      string
	PoetryId     int64
	DynastyName  string
	DynastyUrl   string
	DynastyId    int64
	CategoryList []*TextHrefFormat
	PoetryContentDetail
	PoetryAuthorDetail
}

//诗词详情注释，介绍相关
type PoetryContentDetail struct {
	ContentHtml        string //.sons，所有内容保存在一起， 不同的诗词内容不一样，所以没法单独区分
	Notes              string //注释
	Appreciation       string //赏析
	CreativeBackground string //创作背景
}

//作者信息
type PoetryAuthorDetail struct {
	AuthorName        string
	AuthorId          int64
	AuthorSrcUrl      string
	AuthorTotalPoetry int //诗词总数
}

//诗词对应分类信息
type PoetryCategory struct {
	poetryId    int64
	CategoryIds []int64
}
