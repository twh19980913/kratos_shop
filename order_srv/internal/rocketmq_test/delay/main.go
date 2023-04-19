package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main(){
	p,err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.16.128:9876"}))
	if err != nil {
		fmt.Println("生成producer失败")
		return
	}

	if err = p.Start();err != nil{
		fmt.Println("启动producer失败")
		return
	}

	msg := primitive.NewMessage("imooc1",[]byte("this is delay imooc1"))
	msg.WithDelayTimeLevel(3)

	res,err := p.SendSync(context.Background(),msg)
	if err != nil {
		fmt.Printf("发送失败...%s\n",err)
		return
	}else {
		fmt.Printf("发送成功:%s\n",res.String())
	}

	if err = p.Shutdown();err != nil{
		fmt.Println("关闭producer失败")
		return
	}
}