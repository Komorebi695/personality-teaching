package utils

// RespMsg : 响应数据结构
type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewRespMsg : 生成response对象
func NewRespMsg(code int, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  ResCode(code).StatusText(),
		Data: data,
	}
}
