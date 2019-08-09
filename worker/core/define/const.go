package define

const TestEnv = "test" //单元测试标识

type RedisTaskMark string

const (
	GrabPoetryAll RedisTaskMark = "poetryAll" //抓取诗词所有数据
)

type TaskStatus int

const (
	TaskStatusImplemented TaskStatus = 0  //未执行
	TaskStatusSuccess     TaskStatus = 1  //执行完成
	TaskStatusFail        TaskStatus = -1 //执行失败
)

type RedisKey string

const (
	RedisIsTaskAllRun RedisKey = "task_run_all"
)

type DataFormat string

const (
	HomePoetryCategoryFormatSign DataFormat = "indexPoetryCategory" //首页-诗文分类
	HomePoetryFamousFormatSign   DataFormat = "indexPoetryFamous"   //首页-名句
	HomePoetryAuthorFormatSign   DataFormat = "indexPoetryAuthor"   //首页-作者
)
