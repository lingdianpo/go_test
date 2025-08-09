package initialize

import "go.uber.org/zap"

func LoggerInit() {
	logger, _ := zap.NewProduction()
	//logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
