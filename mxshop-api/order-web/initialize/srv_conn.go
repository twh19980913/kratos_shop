package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/proto"
)

func InitSrvConn()  {
	consulInfo := global.ServerConfig.ConsulInfo
	goodsConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait?wait=14s&tag=srv",consulInfo.Host,consulInfo.Port,global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接 【商品服务失败】")
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	orderConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait?wait=14s&tag=srv",consulInfo.Host,consulInfo.Port,global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接 【订单服务失败】")
	}
	global.OrderSrvClient = proto.NewOrderClient(orderConn)

	invConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait?wait=14s&tag=srv",consulInfo.Host,consulInfo.Port,global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接 【库存服务失败】")
	}
	global.InventorySrvClient = proto.NewInventoryClient(invConn)
}

