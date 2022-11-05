package utils

import (
	"encoding/json"
	"personality-teaching/src/logger"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/bwmarrin/snowflake"
)

const nodeID int64 = 3 // 大于0且小于int64，否则报错
const low int64 = 15   // 切片截取下限
const high int64 = 18  // 切片截取上限

// GenSnowID 生成ID时会上锁，确保不重复
func GenSnowID() string {
	node, _ := snowflake.NewNode(nodeID)
	return node.Generate().String()
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

// SplitNum 截取questionId作为选项与题目的分隔字段
func SplitNum(s string) string {
	return s[low:high]
}

// Obj2Json JSON转换
func Obj2Json(s interface{}) string {
	marshal, _ := json.Marshal(s)
	return string(marshal)
}
