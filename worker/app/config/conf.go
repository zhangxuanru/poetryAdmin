package config

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"poetryAdmin/worker/app/tools"
)

type Config struct {
	RedisHost         string `json:"redis_host"`
	PubChannelTitle   string `json:"pubChannelTitle"`
	GuShiWenIndexUrl  string `json:"gushiwenIndexUrl"`
	GuShiWenPoetryUrl string `json:"gushiwenPoetryUrl"`
	GushiwenSongUrl   string `json:"gushiwenSongUrl"`
	ShiCiMingJuUrl    string `json:"shicimingjuUrl"`
	Env               string `json:"env"`
	DataSource        string `json:"dataSource"`
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
