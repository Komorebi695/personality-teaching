package code

import "github.com/gin-gonic/gin"

// RespMsg : 响应数据结构
type RespMsg struct {
	Code ResCode     `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewRespMsg : 生成response对象
func NewRespMsg(code ResCode, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  code.StatusText(),
		Data: data,
	}
}

func CommonResp(c *gin.Context, httpCode int, statusCode ResCode, data interface{}) {
	c.JSON(httpCode, *NewRespMsg(statusCode, data))
	c.Abort() //此路由后的 gin.HandlerFunc 将不再被调用
}
