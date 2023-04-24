package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/logic"
)

func GetKeyValue(c *gin.Context) {
	respond, _ := logic.GetTencentCloudCOSTemporaryCredentials()
	code.CommonResp(c, http.StatusOK, code.Success, respond)
}
