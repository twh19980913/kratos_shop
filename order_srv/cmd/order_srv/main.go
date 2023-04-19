package main

import (
	"flag"
	"os"
	"encoding/json"
	"fmt"
	"order_srv/internal/conf"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	// "github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"order_srv/internal/pkg"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

type NacosInfo struct {
	Host      string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port      int64  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Namespace string `protobuf:"bytes,3,opt,name=namespace,proto3" json:"namespace,omitempty"`
	User      string `protobuf:"bytes,4,opt,name=user,proto3" json:"user,omitempty"`
	Password  string `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	Dataid    string `protobuf:"bytes,6,opt,name=dataid,proto3" json:"dataid,omitempty"`
}

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	// logger := log.With(log.NewStdLogger(os.Stdout),
	// 	"ts", log.DefaultTimestamp,
	// 	"caller", log.DefaultCaller,
	// 	"service.id", id,
	// 	"service.name", Name,
	// 	"service.version", Version,
	// 	"trace.id", tracing.TraceID(),
	// 	"span.id", tracing.SpanID(),
	// )
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	nacosConfig := &NacosInfo{}
	if err := c.Scan(nacosConfig); err != nil {
		panic(err)
	}

	sc := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.Host, // Nacos的服务地址
			Port:   uint64(nacosConfig.Port),        // Nacos的服务端口
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nacosConfig.Namespace,     // ACM的命名空间Id 当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,                                 // 请求Nacos服务端的超时时间，默认是10000ms
		NotLoadCacheAtStart: true,                                 // 在启动的时候不读取缓存在CacheDir的service信息
		LogDir:              "tmp/nacos/log",   // 日志存储路径
		CacheDir:            "tmp/nacos/cache", // 缓存service信息的目录，默认是当前运行目录
		//RotateTime:          "1h",                                 // 日志轮转周期，比如：30m, 1h, 24h, 默认是24h
		//MaxAge:              3,                                    // 日志最大文件数，默认3
		LogLevel:            "debug",                              // 日志默认级别，值必须是：debug,info,warn,error，默认值是info
	}

	configClient,err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs":sc,
		"clientConfig":cc,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	content,err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.Dataid,
		Group: "dev",
	})
	if err != nil {
		fmt.Println(err)
		return
	}


	var bc conf.Bootstrap
	json.Unmarshal([]byte(content),&bc)

	app, cleanup, err := wireApp(bc.Server, bc.Data,  pkg.Logger())
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
