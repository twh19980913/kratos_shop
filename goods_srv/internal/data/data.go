package data

import (
	"goods_srv/internal/conf"

	"fmt"
	newLog "log"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gorm.io/driver/mysql"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"context"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewNacosClient, NewClient, NewDB,NewBrandsRepo,
	NewBannerRepo,NewCategoryRepo,NewGoodsCategoryBrandRepoRepo,NewGoodsRepo,NewEsClient,NewESClient)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

type ESClient struct{
	esClient *elastic.Client
}

func NewEsClient(c *conf.Data) (*elastic.Client,error){
	//初始化连接
	host := fmt.Sprintf("http://%s:%d",
		c.EsConfig.Host,c.EsConfig.Port)

	logger := newLog.New(os.Stdout, "mxshop", newLog.LstdFlags)
	var err error
	esClient, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	if err != nil {
		log.Fatal(err)
		return nil,err
	}

	//新建mapping 和 index
	//先查询索引是否存在 如果存在就不创建
	exists,err := esClient.IndexExists(EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	if !exists {
		//新建mapping
		_,err = esClient.CreateIndex(EsGoods{}.GetIndexName()).BodyString(EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			log.Fatal(err)
			return nil,err
		}
	}
	return esClient,nil
}

func NewESClient(esClient *elastic.Client) *ESClient {
	return &ESClient{
		esClient: esClient,
	}
}

type NacosClient struct {
	Client *naming_client.INamingClient
}

func NewNacosClient(c *conf.Data, logger log.Logger, client *naming_client.INamingClient) (*NacosClient, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &NacosClient{Client: client}, cleanup, nil
}

func NewClient(c *conf.Data) *naming_client.INamingClient {
	sc := []constant.ServerConfig{
		{
			IpAddr: c.NacosConfig.Host,         // Nacos的服务地址
			Port:   uint64(c.NacosConfig.Port), // Nacos的服务端口
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         c.NacosConfig.Namespace, // ACM的命名空间Id 当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,                    // 请求Nacos服务端的超时时间，默认是10000ms
		NotLoadCacheAtStart: true,                    // 在启动的时候不读取缓存在CacheDir的service信息
		LogDir:              "/tmp/nacos/log",        // 日志存储路径
		CacheDir:            "/tmp/nacos/cache",      // 缓存service信息的目录，默认是当前运行目录
		//RotateTime:          "1h",                                 // 日志轮转周期，比如：30m, 1h, 24h, 默认是24h
		//MaxAge:              3,                                    // 日志最大文件数，默认3
		LogLevel: "debug", // 日志默认级别，值必须是：debug,info,warn,error，默认值是info
	}
	nacosClient, _ := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	return &nacosClient
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}

func NewDB(c *conf.Data) *gorm.DB {
	//dsn := "root:root@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := c.Database.Dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.MysqlConfig.User, c.MysqlConfig.Password,
		c.MysqlConfig.Host, c.MysqlConfig.Port, c.MysqlConfig.Name)
	//
	newLogger := logger.New(
		newLog.New(os.Stdout, "\r\n", newLog.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}
