package main

import (
	"context"
	"fmt"
	v1 "goods_srv/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

var brandsClient v1.GoodsClient

func Init() {
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("127.0.0.1:50052"))
	if err != nil {
		panic(err)
	}

	brandsClient = v1.NewGoodsClient(conn)
}


func TestGetBransList() {
	rsp, err := brandsClient.BrandList(context.Background(), &v1.BrandFilterRequest{

	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

func TestCreateBrand() {
	rsp, err := brandsClient.CreateBrand(context.Background(), &v1.BrandRequest{
		Name: "vivo",
		Logo: "https://edu-991023.oss-cn-beijing.aliyuncs.com/1.jpg",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func main() {
	Init()
	// TestGetUserList()
	// TestCreateUser()
	TestCreateBrand()
}
