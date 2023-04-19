package main

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./myproject.log",
		"stderr",
	}
	return cfg.Build()
}

func main() {
	logger,err := NewLogger()
	if err != nil {
		fmt.Println(err)
		return
	}

	su := logger.Sugar()
	defer su.Sync()
	url := "https://imooc.com"
	su.Info("failed to fetch URL",
		zap.String("url",url),
		zap.Int("attempt",3),
		zap.Duration("backoff",time.Second))
}
