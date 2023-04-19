package main

import (
	"context"
	"fmt"
	v1 "user_srv/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

var userClient v1.UserClient

func Init() {
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("127.0.0.1:9000"))
	if err != nil {
		panic(err)
	}

	userClient = v1.NewUserClient(conn)
}

func TestCreateUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &v1.CreateUserInfo{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1322340383%d", i),
			PassWord: "admin123",
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(rsp.Id)
	}
}

func TestGetUserList() {
	rsp, err := userClient.GetUserLIst(context.Background(), &v1.PageInfo{
		Pn:    1,
		PSize: 2,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)
		checkRsp, err := userClient.CheckPassWord(context.Background(), &v1.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

func main() {
	Init()
	// TestGetUserList()
	TestCreateUser()
}
