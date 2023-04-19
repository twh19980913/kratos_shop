package server

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"syscall"
	v1 "user_srv/api/helloworld/v1"
	"user_srv/internal/conf"
	"user_srv/internal/data"
	"user_srv/internal/service"
	"user_srv/internal/utils"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, userSrv *service.UserService, logger log.Logger,nacosClient *data.NacosClient) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}

	address := strings.Split(c.Grpc.Addr, ":")

	port,_ := utils.GetFreePort()

	realAddr := fmt.Sprintf("%s:%d",address[0],port)

	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(realAddr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterUserServer(srv, userSrv)
	
	ExampleServiceClient_RegisterServiceInstance(*nacosClient.Client, vo.RegisterInstanceParam{
		Ip:          "192.168.16.129",
		Port:        uint64(port),
		ServiceName: "user-srv",
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
			Ip:          "192.168.16.129",
			Port:        uint64(port),
			ServiceName: "user-srv",
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
