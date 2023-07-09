package err_msg

import "errors"

const (
	SUCCESS = 0
	FAILURE = -1 //系统错误码

	St2CustomizeActionResultFailed     = -1 // 失败
	St2CustomizeActionResultSuccess    = 0  // 成功
	St2CustomizeActionResultCodeChange = 1  // 代码变更
	St2CustomizeActionResultEnvCrowded = 5  // 环境挤占
)

/**
 * 系统共用或者需要单独识别的err信息
 */
var (
	DBError    = errors.New("数据库操作失败")
	RedisError = errors.New("redis 操作失败")
	MutexError = errors.New("互斥锁抢占失败")
)
