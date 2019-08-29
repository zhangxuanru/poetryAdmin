package tools

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

//去除html标签
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

//去除正文内容多余的HTML
func TrimDivHtml(str string) (rStr string) {
	mustCompile := regexp.MustCompile(`(?msU)<div class="contyishang">.*</div>`)
	rStr = mustCompile.ReplaceAllString(str, "")
	compile := regexp.MustCompile(`(?msU)<div class="dingpai">.*</div>`)
	rStr = compile.ReplaceAllString(rStr, "")
	nr := regexp.MustCompile(`(?m)[\r\n|\t]`)
	rStr = nr.ReplaceAllString(rStr, "")
	rStr = strings.TrimSpace(rStr)
	return
}

//去除作者诗词列表页，总页数多余文本
//https://so.gushiwen.org/authors/authorvsw_07d17f8539d7A10.aspx
func TrimAuthorTotalPageText(str string) (totalPageNum int, err error) {
	if len(str) == 0 {
		return 0, errors.New("page is nil")
	}
	str = strings.TrimLeft(str, "/")
	str = strings.TrimRight(str, "页")
	str = strings.TrimSpace(str)
	totalPageNum, err = strconv.Atoi(str)
	return
}

//获取URL path
func GetUrlPath(urlStr string) string {
	if strings.Contains(urlStr, "http") == true {
		urlParse, _ := url.Parse(urlStr)
		return urlParse.Path
	}
	return urlStr
}
