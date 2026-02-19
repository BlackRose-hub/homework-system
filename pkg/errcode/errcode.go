package errcode

const (
	Success      = 1000 //操作成功
	InvalidParam = 1001 //请求参数错误
	ServerError  = 1002 //服务器内部错误
	//  1000-1999通用错误
	UserNotFound      = 2001 //用户不存在
	UserAlreadyExists = 2002 //用户已存在
	PassWordIncorrect = 2003 //密码错误
	Unauthorized      = 2004 //未授权
	TokenExpired      = 2005 //token过期
	TokenInvalid      = 2006 //token无效
	PermissionDenied  = 2007 //权限不足
	//2000-2999用户模块错误
	HomeworkNotFound       = 3001 //作业不存在
	HomeworkNotPermission  = 3002 //无权操作该作业
	HomeworkVersionExpired = 3003 //版本过期
	//3000-3999作业模块错误
	SubmissionNotFound       = 4001 //提交记录不存在
	SubmissionAlreadyExists  = 4002 //重复提交
	SubmissionDeadlinePassed = 4003 //超过截至时间
	SubmissionCannotModify   = 4004 //已批改的作业不可修改
	//4000-4999提交模块错误
)

// 错误消息的映射
var Message = map[int]string{
	Success:                  "成功",
	InvalidParam:             "参数错误",
	ServerError:              "服务器内部错误",
	UserNotFound:             "用户不存在",
	UserAlreadyExists:        "用户已存在",
	PassWordIncorrect:        "密码错误",
	Unauthorized:             "未授权",
	TokenExpired:             "token已过期",
	TokenInvalid:             "无效的token",
	PermissionDenied:         "权限不足",
	HomeworkNotFound:         "作业不存在",
	HomeworkNotPermission:    "无权操作该作业",
	HomeworkVersionExpired:   "作业已被修改，请刷新后重试",
	SubmissionNotFound:       "提交记录不存在",
	SubmissionAlreadyExists:  "已提交过该作业",
	SubmissionDeadlinePassed: "已过截至时间",
	SubmissionCannotModify:   "无法修改已批改的作业",
}
