package tools

import (
	"github.com/mozillazg/go-pinyin"
	"strings"
)

//汉字转拼音
func PinYin(han string) string {
	han = strings.TrimSpace(han)
	strs := pinyin.LazyConvert(han, nil)
	if len(strs) > 0 {
		pinYin := strings.Join(strs, "")
		return pinYin
	}
	return ""
}
