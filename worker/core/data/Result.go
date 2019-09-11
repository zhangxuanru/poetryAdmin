package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/core/define"
)

const ChanMaxLen = 50000

//抓取结果处理
type GraspResult struct {
	err       chan error
	close     chan bool
	Data      chan *define.HomeFormat
	ParseData chan *define.ParseData
	storage   *Storage
}

var G_GraspResult *GraspResult

func NewGraspResult() *GraspResult {
	result := &GraspResult{
		err:       make(chan error),
		close:     make(chan bool),
		Data:      make(chan *define.HomeFormat, ChanMaxLen),
		ParseData: make(chan *define.ParseData, ChanMaxLen),
		storage:   NewStorage(),
	}
	G_GraspResult = result
	return result
}

//发送错误消息
func (g *GraspResult) PushError(err error, params ...interface{}) {
	if err != nil {
		go func() {
			g.err <- err
		}()
	}
}

//发送错误消息并关闭协和
func (g *GraspResult) PushErrorAndClose(err error) {
	g.PushError(err)
	g.PushCloseMark(true)
}

//发送结束标志信息
func (g *GraspResult) PushCloseMark(mark bool) {
	g.close <- mark
}

//发送数据
func (g *GraspResult) SendData(data *define.HomeFormat) {
	g.Data <- data
}

func (g *GraspResult) SendParseData(parseData *define.ParseData) {
	g.ParseData <- parseData
}

//统一处理错误消息
func (g *GraspResult) PrintMsg() {
	var (
		err       error
		close     bool
		data      *define.HomeFormat
		parseData *define.ParseData
		autoClose bool
	)
	for {
		if autoClose == true && len(g.ParseData) == 0 {
			goto PRINTERR
		}
		select {
		case err = <-g.err:
			logrus.Debug("Execution error:", err)
			g.WriteErrLog(err)
		case data = <-g.Data:
			g.storage.LoadData(data)
		case parseData = <-g.ParseData:
			parseData.ParseFunc(parseData.Data, parseData.Params)
			logrus.Infoln("g.ParseData len :", len(g.ParseData))
		case close = <-g.close:
			logrus.Infoln("close:", close)
			if len(g.Data) > 0 {
				autoClose = true
				logrus.Info("data 还有数据，暂时不能退出")
				continue
			}
			if close == true {
				goto PRINTERR
			}
		}
	}
PRINTERR:
	logrus.Debug("PrintMsg 结果处理结束......")
	return
}

func (g *GraspResult) WriteErrLog(err error) {
	if err == nil {
		return
	}
	logrus.Infoln("WriteErrLog err:", err)
	//logFile := fmt.Sprintf("error-log-%d-%d-%d.log", time.Now().Year(), time.Now().Month(), time.Now().Hour())
	//file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	//defer file.Close()
	///*
	//		file.Write()
	//		file.WriteString()
	//	 io.WriteString()
	//	ioutil.WriteFile()
	//*/
	//
	//writeObj := bufio.NewWriterSize(file, 4096)
	//buf := []byte(err)
	//if _, err := writeObj.Write(buf); err == nil {
	//	writeObj.Flush()
	//}
}
