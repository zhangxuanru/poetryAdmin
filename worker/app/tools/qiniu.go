package tools

import (
	"bytes"
	"context"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"net/url"
	"strings"
)

//七牛上传图片或MP3
func Upload(src, Bucket, SecretKey, AccessKey string) (fileName string, err error) {
	var (
		byt []byte
	)
	parse, _ := url.Parse(src)
	if len(parse.Path) == 0 {
		return
	}
	urlPath := strings.TrimLeft(parse.Path, "/")
	fileName = strings.ReplaceAll(urlPath, "/", "_")
	if _, byt, err = NewHttpReq().HttpGet(src); err != nil || len(byt) == 0 {
		return
	}
	dataLen := int64(len(byt))
	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = true
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{}

	err = formUploader.Put(context.Background(), &ret, upToken, fileName, bytes.NewReader(byt), dataLen, &putExtra)
	if err != nil {
		return
	}
	return ret.Key, nil
}
