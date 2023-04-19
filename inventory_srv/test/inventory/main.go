package main

import (
	"context"
	"fmt"
	v1 "inventory_srv/api/helloworld/v1"
	"sync"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

var invClient v1.InventoryClient

func Init() {
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("127.0.0.1:50053"))
	if err != nil {
		panic(err)
	}

	invClient = v1.NewInventoryClient(conn)
}

func TestSetInv(goodsId,Num int32) {
	_, err := invClient.SetInv(context.Background(), &v1.GoodsInvInfo{
		GoodsId: goodsId,
		Num: Num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}

func TestInvDetail(goodsId int32) {
	rsp, err := invClient.InvDetail(context.Background(), &v1.GoodsInvInfo{
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Num)
}

func TestSell(wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := invClient.Sell(context.Background(), &v1.SellInfo{
		GoodsInfo: []*v1.GoodsInvInfo{
			{
				GoodsId: 1,
				Num: 1,
			},
			
		},
	})
	if err != nil {
		panic(err)
	}
	
}

func TestReback() {
	_, err := invClient.Reback(context.Background(), &v1.SellInfo{
		GoodsInfo: []*v1.GoodsInvInfo{
			{
				GoodsId: 1,
				Num: 20,
			},
			
		},
	})
	if err != nil {
		panic(err)
	}
	
}

func main() {
	Init()
	// TestGetUserList()
	// TestCreateUser()
	// TestGoodsList()
	// TestGetBatchGoods()
	// TestGoodsDetail()
	// TestCreateGoods()
	// TestInvDetail(1)
	// 并发情况之下，库存无法正确的扣减
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0;i < 10;i++{
		go TestSell(&wg)
	}
	wg.Wait()
	// TestReback()
}