package middle

import (
	"github.com/gin-gonic/gin"
)

const SessionKey = "session_key"

func VerifyTeacherPermission(c *gin.Context) {
	//key, err := c.Cookie(SessionKey)
	//if err == http.ErrNoCookie {
	//	code.CommonResp(c, http.StatusOK, code.NeedLogin, code.EmptyData)
	//}
}
