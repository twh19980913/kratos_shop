package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"order_srv/internal/data"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main2() {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode("generic password", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newPassword)

	passwordInfo := strings.Split(newPassword, "$")
	fmt.Println(passwordInfo)

	check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	fmt.Println(check)
}

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/mxshop_order_srv?charset=utf8mb4&parseTime=True&loc=Local"

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
	_ = db.AutoMigrate(&data.OrderGoods{},&data.OrderInfo{},&data.ShoppingCart{})

	// options := &password.Options{SaltLen: 16,Iterations: 100,KeyLen: 32,HashFunction: sha512.New}
	// salt,encodedPwd := password.Encode("admin123",options)
	// newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)
	// fmt.Println(newPassword)

	// for i := 0;i < 10;i++{
	// 	user := data.User{
	// 		NickName: fmt.Sprintf("bobby%d",i),
	// 		Mobile: fmt.Sprintf("1322340383%d",i) ,
	// 		Password: newPassword,
	// 	}
	// 	db.Save(&user)
	// }
}
