package logic

type Resp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func NewRespData(data interface{}, msg string, code int) *Resp {
	if data == nil {
		data = make(map[string]interface{})
	}
	return &Resp{
		Data: data,
		Msg:  msg,
		Code: code,
	}
}

func (b *Resp) SetCode(code int) {
	b.Code = code
}

func (b *Resp) SetMsg(msg string) {
	b.Msg = msg
}
