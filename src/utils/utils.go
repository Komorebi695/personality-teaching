package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"personality-teaching/src/logger"
	"personality-teaching/src/model"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/bwmarrin/snowflake"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDZWVDfaDbhPNYIU4gUsfawpXXTBQA0xf1nrW+g2pFYED+jDyQk
cumpEl2cAEvF9vJbB7rVIJQFyJfmg0J9XO5X0jDtykJkedHWqi7z9AS056UAXhUQ
cJ+rGwVDu2oBMT/tbCCbRDzuaLcrd5PPQCI1fIrsNQ511cWH6Hv3Lg3JcwIDAQAB
AoGAJtlCDUyRUp0PHJnhnuFYWKaacsdYDBa/foKPi07F39m3piuUqDcp8KBpvvKG
mLHVC9RL3sBd9NKv4/HeNo4fw5pPBjGbqcekYKb2QcA+Mc9HuZojdE6awl1Rf2it
9jCadijlh7LN4U8748E2qXwFnRThSPypQ+mbZWyoLn4EaIECQQD3yku9AmP9qMtH
yk4q/qMx3+Rtq49ohsngSa8EhtrYNjaFBRll7wlFg0p/W9JGb5+IDi4xEmJhTMFl
BfWoFBUhAkEA4IzSuqTwHArKtds2ficGLriiChbt3T1a6XArBg6TsZ7f3x2XoLp2
1j/rsqA3qSmqv9q32pZaowlL6QxAENE4EwJAfkZzbnDnb/8zCPTJ/RMjK2mDyXfi
b0wxWMF0FYR7xi9qfUNp/A5i1S/hKSIr+IUt8XH4jD1oMVmiPM9arzr8wQJAXQm8
Hl1MpzHJf8QONgLRSvZxHSEW+T38txAko2PSyht7wqQuOQhJSMg/TkmYBl0fRFLJ
LqZxc2/cpfjPaqhlRQJAOiO+BzClFXLlhlozxMc1zvVVvFOMXGSlzxdINC7jDbhG
chT5nlMu6vn5sgA6Vb3TRW0w98EmxOLcQwpVNiaD5g==
-----END RSA PRIVATE KEY-----`

const publicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZWVDfaDbhPNYIU4gUsfawpXXT
BQA0xf1nrW+g2pFYED+jDyQkcumpEl2cAEvF9vJbB7rVIJQFyJfmg0J9XO5X0jDt
ykJkedHWqi7z9AS056UAXhUQcJ+rGwVDu2oBMT/tbCCbRDzuaLcrd5PPQCI1fIrs
NQ511cWH6Hv3Lg3JcwIDAQAB
-----END PUBLIC KEY-----`
const (
	nodeID          int64 = 3  // 大于0且小于int64，否则报错
	low             int64 = 15 // 切片截取下限
	high            int64 = 18 // 切片截取上限
	EmptyClassID          = "0"
	DefaultPassWord       = "123456"
	SessionKey            = "session_key"
	TeacherID             = "teacher"
	StudentID             = "student"
	Role                  = "role"
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

// RsaDecrypt 解密前端的RSA加密的密码，返回密码明文
func RsaDecrypt(pwdReq string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(pwdReq)
	if err != nil {
		return []byte{}, err
	}
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error!")
	}
	PK, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Parse private key  error:%s", err))
	}
	return rsa.DecryptPKCS1v15(rand.Reader, PK, b)
}

// CurrentTime 获取当前时间 格式化：2006-01-02 15:01:05
func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:01:05")
}

// SplitNum 截取questionId作为选项与题目的分隔字段
func SplitNum(s string) (string, error) {
	if int64(len(s)) < low || int64(len(s)) < high {
		return "", errors.New("len mismatch")
	}
	return s[low:high], nil
}

// SplitContext 根据题目的分隔字段 分离题干和选项
func SplitContext(id, s string) ([]string, error) {
	splitNum, err := SplitNum(id)
	if err != nil {
		return nil, err
	}
	return strings.Split(s, splitNum), nil
}

// Obj2Json JSON转换
func Obj2Json(s interface{}) (string, error) {
	marshal, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func StuScoreAverage(stu []model.StudentQuestion) map[string]float32 {
	if len(stu) == 0 {
		return nil
	}
	div := make(map[string]float32, len(stu))
	dived := make(map[string]float32, len(stu))
	//根据知识点 求 每个知识点编号的成绩均值，以表示知识点的掌握情况。
	for i := range stu {
		_, ok := dived[stu[i].KnpID]
		if ok {
			dived[stu[i].KnpID] += stu[i].Score / stu[i].AllScore
			div[stu[i].KnpID] += 1
		} else {
			dived[stu[i].KnpID] = stu[i].Score / stu[i].AllScore
			div[stu[i].KnpID] = 1
		}
	}
	for k, _ := range dived {
		dived[k] /= div[k]
	}
	return dived
}
