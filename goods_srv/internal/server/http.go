package server

import (
	//v1 "goods_srv/api/helloworld/v1"
	"goods_srv/internal/conf"
	"goods_srv/internal/service"
	"fmt"
	"strings"
	"goods_srv/internal/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),//日志中间件
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	address := strings.Split(c.Http.Addr, ":")

	port,_ := utils.GetFreePort()

	realAddr := fmt.Sprintf("%s:%d",address[0],port)
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(realAddr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	//v1.RegisterGoodsH(srv, greeter)
	return srv
}
