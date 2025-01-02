package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"happy-dog/utils"
	"time"
)

var DB *gorm.DB
var err error

func InitDb() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", utils.Dbuser, utils.Dbpass, utils.Dbhost, utils.Dbport, utils.Dbname)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("连接数据库失败")
	}

	DB.AutoMigrate(&Customer{}, &Shop{}, &Product{}, &Order{}, &OrderItem{}, &Wallet{}, &Manager{}, &Friends{}, &HistoryMessage{})
	fmt.Printf("数据库迁移成功")

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Printf("获取数据库对象失败")
		return
	}

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	//sqlDB.Close()
}
