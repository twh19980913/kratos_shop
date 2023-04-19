package main

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"sync"
	"time"
)
func main() {
	//创建client
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "127.0.0.1:6379",
	})
	//新建redis连接池
	pool := goredis.NewPool(client)
	//
	rs := redsync.New(pool)

	gNum := 2
	mutexname := "421"

	var wg sync.WaitGroup
	wg.Add(gNum)
	for i := 0; i < gNum; i++ {
		go func() {
			defer wg.Done()
			mutex := rs.NewMutex(mutexname)
			fmt.Println("开始获取锁")
			if err := mutex.Lock(); err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("获取锁成功")
			time.Sleep(time.Second * 5)
			fmt.Println("开始释放锁")
			if ok, err := mutex.Unlock(); !ok || err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("释放锁成功")
		}()
	}
	wg.Wait()

	//var wg sync.WaitGroup
	//wg.Add(gNum)

	//for i := 0; i < gNum; i++ {
	//	go func() {
	//		defer wg.Done()
	//		mutex := rs.NewMutex(mutexname)
	//		fmt.Println("开始获取锁")
	//		if err := mutex.Lock();err != nil{
	//			log.Fatalln(err)
	//		}
	//		fmt.Println("获取锁成功")
	//		time.Sleep(time.Second * 5)
	//		fmt.Println("开始释放锁")
	//		if ok,err := mutex.Unlock(); !ok || err != nil{
	//			log.Fatalln("unlock failed")
	//		}
	//		fmt.Println("释放锁成功")
	//	}()
	//}

	//wg.Wait()
}