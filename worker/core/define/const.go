package define

const TestEnv = "test" //单元测试标识

//redis 订阅发布中发送的任务标识
type RedisTaskMark string

const (
	GrabPoetryAll RedisTaskMark = "poetryAll" //抓取诗词所有数据
)

//任务执行状态
type TaskStatus int

const (
	TaskStatusImplemented TaskStatus = 0  //未执行
	TaskStatusSuccess     TaskStatus = 1  //执行完成
	TaskStatusFail        TaskStatus = -1 //执行失败
)

//redis中的KEY
type RedisKey string

const (
	RedisIsTaskAllRun RedisKey = "task_run_all"
)

//抓取完的数据对应标识
type DataFormat string

const (
	HomePoetryCategoryFormatSign       DataFormat = "indexPoetryCategory"      //首页-诗文分类
	HomePoetryFamousFormatSign         DataFormat = "indexPoetryFamous"        //首页-名句
	HomePoetryAuthorFormatSign         DataFormat = "indexPoetryAuthor"        //首页-作者
	CategoryPoetryAuthorListFormatSign DataFormat = "categoryPoetryAuthorList" //诗文分类-体裁与诗，作者对应列表
)

//首页抓取保存的数据列表
type DataMap map[interface{}]*TextHrefFormat

//诗文分类保存的数据列表
type PoetryDataMap map[interface{}][]interface{}

//数据解析结构
type ParseData struct {
	Data      interface{}
	ParseFunc func(interface{}, interface{})
	Params    interface{}
	DataType  interface{}
}

//显示的位置
type ShowPosition int

const (
	CategoryPosition ShowPosition = 1 //poetry_category表 show_position  1:诗文，2:名句
	FamousPosition   ShowPosition = 2 //poetry_category表 show_position  1:诗文，2:名句
)

//诗词信息内容分类
type DetailNotesType int

const (
	AuthorType DetailNotesType = 2
	PoetryType DetailNotesType = 1
)
