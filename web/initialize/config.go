package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go_test/web/global"
)

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func ConfigInit() {
	//debug := GetEnvInfo("MX_SHOP")
	debug := true
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("web/%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	//文件的路径如何设置
	//v.SetConfigName("config")
	//v.SetConfigType("yaml")
	//v.AddConfigPath("test01/viper_test/")
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//对象如何在其他地方使用
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %v", global.ServerConfig)
	//fmt.Println(v.Get("name"))

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件产生变化: %s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息: %v", global.ServerConfig)
	})
}
