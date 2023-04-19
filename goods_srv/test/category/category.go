package main

import (
	"context"
	"fmt"
	v1 "goods_srv/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var brandsClient v1.GoodsClient

func Init() {
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("127.0.0.1:50052"))
	if err != nil {
		panic(err)
	}

	brandsClient = v1.NewGoodsClient(conn)
}


func TestGetCategoryList() {
	rsp, err := brandsClient.GetAllCategorysList(context.Background(),&emptypb.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.JsonData)
}

func TestGetSubCategoryList() {
	rsp, err := brandsClient.GetSubCategory(context.Background(), &v1.CategoryListRequest{
		Id: 1,

	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.SubCategorys)
}

func main() {
	Init()
	// TestGetUserList()
	// TestCreateUser()
	TestGetSubCategoryList()
}
