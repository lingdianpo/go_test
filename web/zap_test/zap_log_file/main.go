package main

import (
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction() //生产环境
	defer logger.Sync()              // flushes buffer, if any
	url := "https://imooc.com"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
	)

}
