package biz_err

const (
	Success     = 200
	UnKnowErr   = 10000
	ServerError = 10001
	BadRequest  = 10002
)

var CodeMsg = map[int]string{
	Success:     "请求成功",
	UnKnowErr:   "未知业务异常",
	ServerError: "服务端异常",
	BadRequest:  "错误请求",
}

func GetMessage(code int) string {
	if msg, ok := CodeMsg[code]; ok {
		return msg
	}
	return CodeMsg[UnKnowErr]
}
