package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"mxshop-api/goods-web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("goods-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("goods-web/%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	//这个对象如何再其他文件中使用
	//global.ServerConfig = config.ServerConfig{}
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		fmt.Println(err)
		return
	}
	zap.S().Infof("配置信息:%v",global.ServerConfig)
	//fmt.Println(v.Get("name"))

	////viper的功能 - 动态监控变化
	//v.WatchConfig()
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	zap.S().Infof("配置文件产生变化:%v",e.Name)
	//	_ = v.ReadInConfig()
	//	_ = v.Unmarshal(global.ServerConfig)
	//	zap.S().Infof("配置信息:%v",global.ServerConfig)
	//})

	//从服务器nacos读取配置
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port: global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig {
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		zap.S().Fatalln(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		zap.S().Fatalln(err)
	}
	//fmt.Println(content) //字符串 - yaml
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil{
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}
	fmt.Println(&global.ServerConfig)
}
