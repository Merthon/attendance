package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 配置你的数据库账号密码
	// 格式: 用户名:密码@tcp(IP:端口)/数据库名?参数
	dsn := "xxxxx:xxxxxx@tcp(127.0.0.1:3306)/attendance_db?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	fmt.Println("数据库连接成功！")
    
    // 自动迁移
	// database.AutoMigrate(&model.User{}) 
}