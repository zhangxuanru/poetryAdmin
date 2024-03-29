package tools

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	ctx := md5.New()
	ctx.Write([]byte(str))
	return hex.EncodeToString(ctx.Sum(nil))
}
