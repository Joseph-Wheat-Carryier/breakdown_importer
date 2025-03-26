package gorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	connectionString := "root:KF@32rjb@tcp(172.16.1.62:3306)/cnxm?charset=utf8mb4&parseTime=True&loc=Local"
	initDB(connectionString)
}

var db *gorm.DB

func initDB(connectionString string) error {
	// 使用数据库驱动初始化数据库连接
	connection, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}

	// 将数据库连接赋给全局变量
	db = connection

	// 迁移数据库表（如果需要的话）
	// db.AutoMigrate(&YourModel{})

	return nil
}

func GetDB() *gorm.DB {
	return db
}
