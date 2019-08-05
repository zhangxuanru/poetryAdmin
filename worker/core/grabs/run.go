package grabs

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/redis"
	"poetryAdmin/worker/core/parse"
)

func Run() (err error) {
	var (
		subReceiveChan chan []byte
	)
	subReceiveChan = make(chan []byte)
	go parse.NewAnalysis(subReceiveChan).ParseSubscribeData()
	if err = redis.SubScribe(config.G_Conf.PubChannelTitle, subReceiveChan); err != nil {
		logrus.Debug("SubScribe err:", err)
		return err
	}
	return nil
}
