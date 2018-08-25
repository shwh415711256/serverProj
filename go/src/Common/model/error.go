package model

type ErrorType int

const (
	InnerError ErrorType = 1				// 内部错误
	RequestParamError ErrorType = 2		// 参数错误
)

type ErrorStatus int
const (
	LoginParamErrorStatus ErrorStatus	= 100			// 登录参数出错
	WXLoginErrorStatus ErrorStatus = 101				// 微信登陆出错
	WXUpdateErrorStatus ErrorStatus = 102				// 同步微信用户数据出错
	CreateRoomParamErrorStatus ErrorStatus = 103		// 创建房间参数出错
	GetOneMatchInfoRpcErrorStatus ErrorStatus = 104		// 获取一个匹配队友信息出错
	CommitGameInfoErrorStatus ErrorStatus = 105			// 上传游戏历史数据出错
	CommitGameInfoRpcErrorStatus ErrorStatus = 106		// 上传游戏数据rpc调用出错
	EnrollMatchErrorStatus ErrorStatus = 107				// 报名比赛出错
	GetMatchTicketRpcErrorStatus ErrorStatus = 108		// 获取比赛门票信息出错
)

type CommonFormatError struct {
	Type ErrorType
	Status ErrorStatus
	ErrorDesc string
}
