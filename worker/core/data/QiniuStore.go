package data

import (
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
)

type UploadStore struct {
}

func NewUploadStore() *UploadStore {
	return &UploadStore{}
}

//七牛上传图片或MP3
func (i *UploadStore) Upload(src string) (fileName string, err error) {
	if len(src) == 0 {
		return
	}
	fileName, err = tools.Upload(src, config.G_Conf.Bucket, config.G_Conf.SecretKey, config.G_Conf.AccessKey)
	return
}
