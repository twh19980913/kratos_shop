package server

import (
	"fmt"
	v1 "goods_srv/api/helloworld/v1"
	"goods_srv/internal/conf"
	"goods_srv/internal/data"
	"goods_srv/internal/service"
	// "goods_srv/internal/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GoodsService, logger log.Logger,nacosClient *data.NacosClient) *grpc.Server {
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

	port := 50052

	realAddr := fmt.Sprintf("%s:%d",address[0],port)

	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(realAddr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterGoodsServer(srv, greeter)
	ExampleServiceClient_RegisterServiceInstance(*nacosClient.Client, vo.RegisterInstanceParam{
		Ip:          c.Host,
		Port:        uint64(port),
		ServiceName: c.Name,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	go func() {
		// acquire stop sign
		quit := make(chan os.Signal)
		signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
		<- quit
		ExampleServiceClient_DeRegisterServiceInstance(*nacosClient.Client, vo.DeregisterInstanceParam{
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