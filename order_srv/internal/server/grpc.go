package server

import (
	v1 "order_srv/api/helloworld/v1"
	"order_srv/internal/conf"
	"order_srv/internal/service"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"fmt"
	"order_srv/internal/data"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, order *service.OrderService, logger log.Logger,nacosClient *data.NacosClient) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),//日志中间件
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	address := strings.Split(c.Grpc.Addr, ":")

	port := 50054

	realAddr := fmt.Sprintf("%s:%d",address[0],port)
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(realAddr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterOrderServer(srv, order)
	ExampleServiceClient_RegisterServiceInstance(nacosClient.Client, vo.RegisterInstanceParam{
		Ip:          c.Host,
		Port:        uint64(port),
		ServiceName: c.Name,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	// 这里 订阅消息 订单超时归还topic

	rockClient,err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.16.128:9876"}),
		consumer.WithGroupName("mxshop-order"),
	)
	if err != nil {
		fmt.Println("生成consumer失败")
	}

	if err = rockClient.Subscribe("order_timeout",consumer.MessageSelector{},func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		return order.OrderTimeout(ctx,me...)
	});err != nil{
		fmt.Println("读取消息失败")
	}

	_ = rockClient.Start()

// 监听结束

	go func() {
		// acquire stop sign
		quit := make(chan os.Signal)
		signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
		<- quit
		ExampleServiceClient_DeRegisterServiceInstance(nacosClient.Client, vo.DeregisterInstanceParam{
			Ip:          c.Host,
			Port:        uint64(port),
			ServiceName: c.Name,
			Ephemeral:   true, //it must be true
		})
	}()
	return srv
}
func ExampleServiceClient_DeRegisterServiceInstance(client naming_client.INamingClient, param vo.DeregisterInstanceParam) {
	success, _ := client.DeregisterInstance(param)
	fmt.Printf("注销服务实例:%+v,result:%+v \n\n", param, success)
}


func ExampleServiceClient_RegisterServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) {
	success, _ := client.RegisterInstance(param)
	fmt.Printf("注册服务实例:%+v,result:%+v \n\n", param, success)
}