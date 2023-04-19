package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"time"
)

func main(){
	// c,err := rocketmq.NewPushConsumer(
	// 	consumer.WithNameServer([]string{"192.168.16.128:9876"}),
	// 	consumer.WithGroupName("mxshop"),
	// )
	// if err != nil {
	// 	fmt.Println("生成consumer失败")
	// 	return
	// }

	// if err = c.Subscribe("imooc1",consumer.MessageSelector{},func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// 	for i := range msgs{
	// 		fmt.Printf("获取到值: %v\n",msgs[i])
	// 	}
	// 	return consumer.ConsumeSuccess,nil
	// });err != nil{
	// 	fmt.Println("读取消息失败")
	// }

	// _ = c.Start()

	// //不能让主goroutine退出
	// time.Sleep(time.Hour)
	// _ = c.Shutdown()

	rockClient,err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.16.128:9876"}),
		consumer.WithGroupName("mxshop-inventory"),
	)
	if err != nil {
		fmt.Println("生成consumer失败")
	}

	if err = rockClient.Subscribe("order_reback",consumer.MessageSelector{},func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgs{
			fmt.Printf("获取到值: %v\n",msgs[i])
		}
		return consumer.ConsumeSuccess,nil
	});err != nil{
		fmt.Println("读取消息失败")
	}

	_ = rockClient.Start()
	//不能让主goroutine退出
	time.Sleep(time.Hour)
	_ = rockClient.Shutdown()
}