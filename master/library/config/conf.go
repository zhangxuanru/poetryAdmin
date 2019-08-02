package config

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"poetryAdmin/master/library/tools"
)

type Config struct {
	ApiPort         int    `json:"apiPort"`
	ApiReadTimeOut  int    `json:"apiReadTimeOut"`
	ApiWriteTimeOut int    `json:"apiWriteTimeOut"`
	WebRoot         string `json:"webroot"`
	ViewDir         string `json:"viewDir"`
	StaticDomain    string `json:"staticDomain"`
	BaseDomain      string `json:"baseDomain"`
}

var (
	G_Json jsoniter.API
	G_Conf *Config
)

func InitConfig(file string) (err error) {
	var (
		conf    Config
		content []byte
		exists  bool
	)
	G_Json = jsoniter.ConfigCompatibleWithStandardLibrary
	if exists, err = tools.PathExists(file); err != nil || exists == false {
		goto PRINTERR
	}
	if content, err = ioutil.ReadFile(file); err != nil {
		goto PRINTERR
	}
	if err = G_Json.Unmarshal(content, &conf); err != nil {
		goto PRINTERR
	}
	G_Conf = &conf
	return nil
PRINTERR:
	return err
}
