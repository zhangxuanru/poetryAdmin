package data

import (
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
)

type ImgStore struct {
}

func NewImgStore() *ImgStore {
	return &ImgStore{}
}

//上传图片
func (i *ImgStore) UploadImg(src string) (fileName string, err error) {
	fileName, err = tools.UploadImg(src, config.G_Conf.Bucket, config.G_Conf.SecretKey, config.G_Conf.AccessKey)
	return
}
