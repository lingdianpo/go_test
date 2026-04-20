package initialize

import "go.uber.org/zap"

func LoggerInit() {
	//logger, _ := zap.NewProduction()
	//logger, _ := zap.NewDevelopment()
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"stdout",
		"./my_project.log",
	}
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}
