package tools

import (
	"hash/crc32"
)

//crc32 加密
func Crc32(str string) uint32 {
	ieee := crc32.ChecksumIEEE([]byte(str))
	return ieee
}
