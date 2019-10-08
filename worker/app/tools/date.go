/*
@Time : 2019/9/16 17:57
@Author : zxr
@File : date
@Software: GoLand
*/
package tools

import (
	"fmt"
	"time"
)

func GetCurrentUnix() int64 {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	now := time.Now().In(cstSh)
	formStr := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), 0, 0, 0)
	parse, _ := time.ParseInLocation("2006-01-02 15:04:05", formStr, cstSh)
	return parse.Unix()
}
