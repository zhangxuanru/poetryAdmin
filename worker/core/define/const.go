package define

type RedisTaskMark string

const (
	GrabPoetryAll RedisTaskMark = "All" //抓取所有数据
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
