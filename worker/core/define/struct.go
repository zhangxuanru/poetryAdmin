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
	Title              string
	Content            string
	PoetryId           int64
	DynastyName        string
	DynastyUrl         string
	DynastyId          int64
	CategoryList       []*TextHrefFormat
	CreativeBackground string //创作背景
	PoetryAuthorDetail
	Detail []*PoetryContentData
}

//----诗词正文具体数据 诗词详情注释，介绍相关
type PoetryContentData struct {
	AppRecId   int    //赏析ID
	TransId    int    //翻译ID
	Sort       int    //排序
	Title      string //标题
	Content    string //具体内容
	PlayUrl    string //源声音路径
	PlaySrcUrl string //源声音API路径
	FileName   string //下载的文件路径
}

//作者信息
type PoetryAuthorDetail struct {
	AuthorName        string
	AuthorId          int64
	AuthorSrcUrl      string //作者介绍页
	AuthorContentUrl  string //作者诗词列表页
	AuthorImgUrl      string
	AuthorTotalPoetry int //诗词总数
}

//诗词对应分类信息
type PoetryCategory struct {
	poetryId    int64
	CategoryIds []int64
}
