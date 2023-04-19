package main

import (
	"context"
	"fmt"
	v1 "goods_srv/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

var goodsClient v1.GoodsClient

func Init() {
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("127.0.0.1:50052"))
	if err != nil {
		panic(err)
	}

	goodsClient = v1.NewGoodsClient(conn)
}

func TestGoodsList() {
	rsp, err := goodsClient.GoodsList(context.Background(), &v1.GoodsFilterRequest{
		TopCategory: 2,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Data)
}

func TestCreateGoods() {
	rsp, err := goodsClient.CreateGoods(context.Background(), &v1.CreateGoodsInfo{
		Name:            "华为P70",
		GoodsSn:         "13273145995",
		Stocks:          20,
		CategoryId:      3,
		BrandId:         2,
		MarketPrice:     99,
		ShopPrice:       99,
		GoodsBrief:      "超强华为手机",
		Images:          []string{"https://edu-991023.oss-cn-beijing.aliyuncs.com/1.jpg"},
		DescImages:      []string{"https://edu-991023.oss-cn-beijing.aliyuncs.com/1.jpg"},
		GoodsFrontImage: "https://edu-991023.oss-cn-beijing.aliyuncs.com/1.jpg",
		ShipFree:        true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func TestGetBatchGoods() {
	rsp, err := goodsClient.BatchGetGoods(context.Background(), &v1.BatchGoodsIdInfo{
		Id: []int32{1},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Data)
}

func TestGoodsDetail() {
	rsp, err := goodsClient.GetGoodsDetail(context.Background(), &v1.GoodInfoRequest{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}

func main() {
	Init()
	// TestGetUserList()
	// TestCreateUser()
	// TestGoodsList()
	// TestGetBatchGoods()
	// TestGoodsDetail()
	// TestCreateGoods()
	TestGoodsDetail()
}
