package utils

import (
	"personality-teaching/src/logger"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/bwmarrin/snowflake"
)

const nodeID int64 = 3 // 大于0且小于int64，否则报错

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

// CurrentTime 获取当前时间 格式化：2006-01-02 15:01:05
func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:01:05")
}
