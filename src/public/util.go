package public

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
)

func GenSaltPassword(salt, password string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}

//MD5 md5加密
func MD5(s string) string {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//JSON转换
func Obj2Json(s interface{}) string {
	marshal, _ := json.Marshal(s)
	return string(marshal)
}

//UUID生成
func UUID() string {
	return uuid.New().String()
}
