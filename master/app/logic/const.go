package logic

const RespSuccess = 200
const RespFail = -1

const LoginCookieName = "poetryAdmin"
const LoginCookieUserName = "UserName"
const LoginCookiePassword = "PassWord"

const (
	RespSuccessMsg   = "请求成功"
	RespFailMsg      = "请求失败"
	RespLoginSuccess = "登录成功"
	RespLoginFailMsg = "登录失败"

	GrabTaskTitleAll = "抓取所有数据"
	GrabTaskAdd      = "添加任务成功,等待执行"
)

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
