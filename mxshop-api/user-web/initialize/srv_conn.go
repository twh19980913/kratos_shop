package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func InitSrvConn()  {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait?wait=14s&tag=srv",consulInfo.Host,consulInfo.Port,global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[GetUserList] 连接 【用户服务失败】")
	}
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

//初始化grpc服务
func InitSrvConnTest() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d",consulInfo.Host, consulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		log.Fatalln(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}
	//拨号连接用户grpc服务
	userConn,err := grpc.Dial(fmt.Sprintf("%s:%d",userSrvHost,userSrvPort),grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg",err.Error())
	}
	//1、后续的用户服务下线了 2、改端口了 3、改IP了
	//2、已经事先创建好了连接，这样后续就不用进行两次tcp的三次握手了
	//3、多个goroutine共用一个连接会造成性能问题，这里我们使用连接池来优化
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
