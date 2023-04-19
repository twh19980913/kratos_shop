package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
)

func InitSrvConn()  {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait?wait=14s&tag=srv",consulInfo.Host,consulInfo.Port,global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[GetUserList] 连接 【用户服务失败】")
	}
	userSrvClient := proto.NewGoodsClient(userConn)
	global.GoodsSrvClient = userSrvClient
}

