package mysql

import (
	"errors"
	"fmt"
	"log"
	"os"
	"personality-teaching/src/configs"
	"personality-teaching/src/model"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func InitMysql(initConfig *configs.AppConfig) (err error) {
	database := initConfig.DataBase

	// 日志设置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别 Silent、Error、Warn、Info
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		database.User,
		database.Pwd,
		database.Host,
		database.Port,
		database.Database)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return
	}

	// 迁移
	err = Db.AutoMigrate(&model.KnowledgePointFile{},&model.QuestionFile{})
	if err != nil {
		return err
	}

	return err
}
func GetGormPool() (*gorm.DB, error) {
	if Db != nil {
		return Db, nil
	} else {
		return nil, errors.New("get pool error")
	}
}
