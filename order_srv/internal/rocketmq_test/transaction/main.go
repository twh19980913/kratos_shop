package main

import (
	"context"
	"fmt"
	"time"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type OrderListener struct {
	
}

func (o *OrderListener) ExecuteLocalTransaction (msg *primitive.Message) primitive.LocalTransactionState {
	// fmt.Println("开始执行本地逻辑")
	// time.Sleep(time.Second * 3)
	// fmt.Println("执行本地逻辑成功")
	// return primitive.CommitMessageState

	//fmt.Println("开始执行本地逻辑")
	//time.Sleep(time.Second * 3)
	//fmt.Println("执行本地逻辑失败")
	//return primitive.RollbackMessageState // 消息并不会写到broker

	fmt.Println("开始执行本地逻辑")
	time.Sleep(time.Second * 3)
	fmt.Println("执行本地逻辑失败")
	//本地执行逻辑无缘无故失败 代码异常 宕机
	return primitive.UnknowState
}

func (o *OrderListener) CheckLocalTransaction (msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("rocketmq的消息回查")
	time.Sleep(time.Second * 15)
	//本地事务执行失败了 消息不要投递出去
	return primitive.CommitMessageState
}

func main(){
	p,err := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNameServer([]string{"192.168.16.128:9876"}),
	)
	if err != nil {
		fmt.Println("生成producer失败")
		return
	}

	if err = p.Start();err != nil{
		fmt.Println("启动producer失败")
		return
	}


	res,err := p.SendMessageInTransaction(context.Background(),primitive.NewMessage("TransTopic",[]byte("this is transaction unknow")))

	if err != nil {
		fmt.Printf("发送失败...%s\n",err)
		return
	}else {
		fmt.Printf("发送成功:%s\n",res.String())
	}

	time.Sleep(time.Hour)

	if err = p.Shutdown();err != nil{
		fmt.Println("关闭producer失败")
		return
	}
}