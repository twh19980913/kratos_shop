package main

import (
	"fmt"
	"inventory_srv/internal/data"
	"log"
	"os"
	"time"
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/mxshop_inventory_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//设置全局的logger,这个logger在我们执行每个sql语句的时候都会打印一行Sql
	//sql才是最重要的，

	//定义一个表结构，将表结构直接生成对应的表-migrations
	//迁移 schema
	// _ = db.AutoMigrate(&data.StockSellDetail{})

	// options := &password.Options{SaltLen: 16,Iterations: 100,KeyLen: 32,HashFunction: sha512.New}
	// salt,encodedPwd := password.Encode("admin123",options)
	// newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)
	// fmt.Println(newPassword)
	// orderDetail := data.StockSellDetail{
	// 	OrderSn: "imooc-bobby",
	// 	Status: 1,
	// 	Detail: []data.GoodsDetail{
	// 		{
	// 			Goods: 1,
	// 			Num: 2,
	// 		},
	// 		{
	// 			Goods: 4,
	// 			Num: 2,
	// 		},
	// 	},
	// }
	// db.Create(&orderDetail)

	var sellDetail data.StockSellDetail
	db.Where(data.StockSellDetail{OrderSn: "imooc-bobby"}).First(&sellDetail)
	fmt.Println(sellDetail.Detail)
	// for i := 0;i < 10;i++{
	// 	user := data.User{
	// 		NickName: fmt.Sprintf("bobby%d",i),
	// 		Mobile: fmt.Sprintf("1322340383%d",i) ,
	// 		Password: newPassword,
	// 	}
	// 	db.Save(&user)
	// }
}