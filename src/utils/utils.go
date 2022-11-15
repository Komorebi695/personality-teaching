package utils

import (
	"encoding/json"
	"personality-teaching/src/logger"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/bwmarrin/snowflake"
)

const (
	nodeID          int64 = 3  // 大于0且小于int64，否则报错
	low             int64 = 15 // 切片截取下限
	high            int64 = 18 // 切片截取上限
	EmptyClassID          = "0"
	DefaultPassWord       = "123456"
	SessionKey            = "session_key"
	TeacherID             = "teacher"
)

// GenSnowID 生成ID时会上锁，确保不重复
func GenSnowID() string {
	node, _ := snowflake.NewNode(nodeID)
	return node.Generate().String()
}

// GetUUID 生成uuid
func GetUUID() string {
	return uuid.New().String()
}

// Encryption 将传入的字符串进行hash加密
func Encryption(toHash string) (string, error) {
	hashedID, err := bcrypt.GenerateFromPassword([]byte(toHash), 0)
	if err != nil {
		logger.L.Error("bcrypt加密失败：", zap.Error(err))
		return "", err
	}
	return string(hashedID), nil
}

// GetDefaultPassWord 返回加密后的默认密码
func GetDefaultPassWord() string {
	s, _ := Encryption(DefaultPassWord)
	return s
}

// CompareHash 请确保传入的hashedStr 不为空字符串，否则将报错
func CompareHash(hashedStr string, str string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		logger.L.Error("bcrypt 比对比失败：", zap.Error(err))
		return false, err
	}
	return true, nil
}

// CurrentTime 获取当前时间 格式化：2006-01-02 15:01:05
func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:01:05")
}

// SplitNum 截取questionId作为选项与题目的分隔字段
func SplitNum(s string) string {
	return s[low:high]
}

// SplitContext 根据题目的分隔字段 分离题干和选项
func SplitContext(id, s string) []string {
	splitNum := SplitNum(id)
	return strings.Split(s, splitNum)
}

// Obj2Json JSON转换
func Obj2Json(s interface{}) string {
	marshal, _ := json.Marshal(s)
	return string(marshal)
}
