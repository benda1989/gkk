package code

// app(rpcs) part model
//
//	0001     0    00
const (
	SUCCESS = 0 // 操作成功

	ERROR           = 10000 // 默认错误返回
	SERVER_ERROR    = 10001 // 服务器错误
	UNKNOW_ERROR    = 10002 // 未知错误
	REMOTE_ERROR    = 10003 // 远程服务错误
	IP_LINMIT       = 10004 // 地址限制
	IP_LINMIT_CODE  = 10014 // 地址限制需要验证码
	NO_PERMISSION   = 10005 // 未拥有授权
	PARAMETER_ERROR = 10008 // 参数错误
	RPC_ERROR       = 10013 // RPC通讯错误
	NOT_FOUND_ROUTE = 10020 // 未查询到路由
	NOT_FOUND_METH  = 10021 // 未查询到方式

	NOT_FOUND      = 10022 // 未查询到
	AUTH_ERROR     = 10023 // 认证错误
	NO_VALID_TOKEN = 10024 // token无效
	REPEAD         = 10025 // 重复数据或操作
	OUT_SLINCE     = 10026 // 超出限制
	DB_ERRROR      = 10027 // 数据库错误
	SENSITIVE      = 10028 // 敏感词语

	SMS_SEND_ERROR = 10029 //验证码发送失败
	SMS_LINMIT     = 10030 //验证码超过限制
	CODE_EXPIRE    = 10031 //验证码过期
	CODE_WRONG     = 10032 //验证码错误
	CODE_SENS      = 10033 //敏感词
)
const (
	//用户
	NOT_FOUND_USER   = 20001 // 未查询到用户
	EXIST_USER       = 20002 // 用户已存在
	NO_BELONG_ACTION = 20003 // 越权操作
	NO_PHONE         = 20004 // 未验证手机号

	//订单
	ORDER_ERROR     = 20100 // 订单统用错误
	NO_ORDER        = 20101 // 订单不存在
	NO_PAID_ORDER   = 20102 // 订单未支付
	PAID_ORDER      = 20103 // 订单已支付
	EXPIRE_ORDER    = 20104 // 订单超时
	DONE_ORDER      = 20105 // 订单已完结
	NO_BELONG_ORDER = 20106 // 不属于自己的订单

	//积分
	SIGNED             = 20200 // 已签到
	NO_ENOUGH_POINT    = 20201 // 积分不足
	NO_ILLEGAL_CHANNEL = 20202 // 渠道非法
	NO_ILLEGAL_SIGN    = 20203 // 渠道非签到

)
