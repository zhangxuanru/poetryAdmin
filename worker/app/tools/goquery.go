package tools

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

//通过文本内容实例化goquery
func NewDocumentFromReader(html string) (query *goquery.Document, err error) {
	return goquery.NewDocumentFromReader(strings.NewReader(html))
}

//通过http 返回值  实例化goquery
func NewDocumentFromResponse(res *http.Response) (query *goquery.Document, err error) {
	return goquery.NewDocumentFromResponse(res)
}
