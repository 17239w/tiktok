package models

import (
	"tiktok/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// 连接数据库
	DB, err = gorm.Open(mysql.Open(config.DBConnectString()), &gorm.Config{
		PrepareStmt:            true, //缓存预编译命令
		SkipDefaultTransaction: true, //禁用默认事务操作
		//Logger:                 logger.Default.LogMode(logger.Info), // 设置日志级别为Info，打印SQL语句
	})
	if err != nil {
		panic(err)
	}
	// 自动迁移
	err = DB.AutoMigrate(&UserInfo{}, &UserLogin{})
	if err != nil {
		panic(err)
	}
}
