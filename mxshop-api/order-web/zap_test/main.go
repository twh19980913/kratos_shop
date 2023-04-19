package main

import (
	"go.uber.org/zap"
)

func main() {
	logger,_ := zap.NewProduction()
	defer logger.Sync()
	url := "http://imooc.com"
	logger.Info("failed to fetch URL",
		zap.String("url",url),
		zap.Int("nums",3))
	//sugar := logger.Sugar()
	//sugar.Infow("failed to fetch URL",
	//	"url",url,
	//	"attempt",3)
	//sugar.Infof("Failed to fetch URL: %s",url)
}
