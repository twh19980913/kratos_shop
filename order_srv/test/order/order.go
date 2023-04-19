package main

import (
	"context"
	
	"fmt"
	v1 "order_srv/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	// "google.golang.org/protobuf/types/known/emptypb"
)

var orderClient v1.OrderClient

func Init() {
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("127.0.0.1:50054"))
	if err != nil {
		panic(err)
	}

	orderClient = v1.NewOrderClient(conn)
}


func TestCreateCartItem(userId,nums,goodsId int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(),&v1.CartItemRequest{
		UserId: userId,
		Nums: nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func TestCartItemList(userId int32){
	rsp, err := orderClient.CartItemList(context.Background(),&v1.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}
	for _,item := range rsp.Data{
		fmt.Println(item.Id,item.GoodsId,item.Nums)
	}
}

func TestUpdateCartItem(id int32){
	_, err := orderClient.UpdateCartItem(context.Background(),&v1.CartItemRequest{
		Id: id,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateOrder(){
	_,err := orderClient.CreateOrder(context.Background(),&v1.OrderRequest{
		UserId: 21,
		Address: "北京市",
		Name: "bobby",
		Mobile: "13223403830",
		Post: "请尽快发货",
	})
	if err != nil {
		panic(err)
	}
}

func TestGetOrderDetail(orderId int32){
	rsp, err := orderClient.OrderDetail(context.Background(),&v1.OrderRequest{
		Id: orderId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.OrderInfo.OrderSn)

	for _,good := range rsp.Goods{
		fmt.Println(good.GoodsName)
	}
}

func TestOrderList(){
	rsp, err := orderClient.OrderList(context.Background(),&v1.OrderFilterRequest{
		
	})
	if err != nil {
		panic(err)
	}
	
	for _,order := range rsp.Data{
		fmt.Println(order.OrderSn)
	}
}

func main() {
	Init()
	// TestGetUserList()
	// TestCreateUser()
	// TestCreateCartItem(21,1,4)
	//TestCartItemList(21)
	// TestUpdateCartItem(1)
	// TestCreateOrder()
	// TestGetOrderDetail(1)
	TestOrderList()
}