package initialize

import "go.uber.org/zap"

func LoggerInit() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}
