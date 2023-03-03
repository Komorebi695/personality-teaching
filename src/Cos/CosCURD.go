package Cos

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/model"
	"strconv"
	"time"
)

// KnpUploadFileToCos 上传知识点文件
func KnpUploadFileToCos(c *gin.Context){
	secretId := "AKIDpUwCFD0ZtwwrZ2GG1mMm8UsQ3dyCF4h5"
	secretKey := "FF0VpuUxEJF1Yq00pVN9jBi182h7JMBk"
	region := "ap-nanjing"
	bucket := "teachlearning-1314366587"

	Client, err := NewCosClient(secretId, secretKey, region, bucket)
	if err != nil {
		fmt.Printf("failed to connect COS: %s\n", err.Error())
		return
	}

    file, err := c.FormFile("file") // 获取上传的文件
	if err != nil{
		return
	}

	f,err := file.Open()
	if err!=nil{
		return
	}
	defer f.Close()
	//对文件名进行加密
	filename := file.Filename // 获取文件名
	ext := filepath.Ext(filename)//获取文件的后缀名
	hasher := sha256.New()
	hasher.Write([]byte(file.Filename))
	encryptedFilename := hex.EncodeToString(hasher.Sum(nil)) //加密文件名
	encryptedFile :=  strconv.Itoa(int(time.Now().Unix()))+"-" + encryptedFilename+ext//加上文件名后缀
	cosPrefix := "https://teachlearning-1314366587.cos.ap-nanjing.myqcloud.com"
	cosFilename := cosPrefix + "/" + strconv.Itoa(int(time.Now().Unix())) + "-" + encryptedFilename + ext


	_, err = Client.Object.Put(context.Background(), encryptedFile, f, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	KnpFileModel := &model.KnowledgePointFile{
		CosUrl: cosFilename,
	}
	result := mysql.Db.Create(&KnpFileModel)
	c.JSON(http.StatusOK,gin.H{
		"errno":0,
		"data":gin.H{
			"url":cosFilename,
			"href":"",
			"alt":"",
		},
		"result":result.Error,
	})
}

// QuestionUploadFileToCos 上传题目文件
func QuestionUploadFileToCos(c *gin.Context){
	secretId := "AKIDpUwCFD0ZtwwrZ2GG1mMm8UsQ3dyCF4h5"
	secretKey := "FF0VpuUxEJF1Yq00pVN9jBi182h7JMBk"
	region := "ap-nanjing"
	bucket := "teachlearning-1314366587"

	Client, err := NewCosClient(secretId, secretKey, region, bucket)
	if err != nil {
		fmt.Printf("failed to connect COS: %s\n", err.Error())
		return
	}

	file, err := c.FormFile("file") // 获取上传的文件
	if err != nil{
		return
	}

	f,err := file.Open()
	if err!=nil{
		return
	}
	defer f.Close()
	//对文件名进行加密
	filename := file.Filename // 获取文件名
	ext := filepath.Ext(filename)//获取文件的后缀名
	hasher := sha256.New()
	hasher.Write([]byte(file.Filename))
	encryptedFilename := hex.EncodeToString(hasher.Sum(nil)) //加密文件名
	encryptedFile :=  strconv.Itoa(int(time.Now().Unix()))+"-" + encryptedFilename+ext//加上文件名后缀
	cosPrefix := "https://teachlearning-1314366587.cos.ap-nanjing.myqcloud.com"
	cosFilename := cosPrefix + "/" + strconv.Itoa(int(time.Now().Unix())) + "-" + encryptedFilename + ext
	//cosFilename := cosPrefix +"/"+file.Filename

	_, err = Client.Object.Put(context.Background(), encryptedFile, f, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	QuestionFileModel := &model.QuestionFile{
		CosUrl: cosFilename,
	}
	result := mysql.Db.Create(&QuestionFileModel)
	c.JSON(http.StatusOK,gin.H{
		"errno":0,
		"data":gin.H{
			"url":cosFilename,
			"href":"",
			"alt":"",
		},
		"result":result.Error,
	})
}



