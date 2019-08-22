package tools

import (
	"hash/crc32"
	"math/rand"
)

//crc32 加密
func Crc32(str string) uint32 {
	ieee := crc32.ChecksumIEEE([]byte(str))
	return ieee
}

//生成区间随机数
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
