package define

//
type HomeFormat struct {
	Identifier DataFormat
	Data       interface{}
}

//首页保存的数据格式
type TextHrefFormat struct {
	Text string
	Href string
}
