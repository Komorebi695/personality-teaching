package utils

type ResCode int

// 业务逻辑状态码
const (
	LoginState ResCode = iota

	// UnKnowError 未知错误
	UnKnowError
)

//业务逻辑状态信息描述
var recodeText = map[ResCode]string{
	LoginState:  "请登录",
	UnKnowError: "error",
}

// StatusText 返回状态码的文本。如果代码为空或未知状态码则返回error
func (code ResCode) StatusText() string {
	msg, ok := recodeText[code]
	if ok {
		return msg
	}
	return recodeText[UnKnowError]
}
