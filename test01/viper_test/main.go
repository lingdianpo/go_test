package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

func GetEnvInfo(env string) int {
	viper.AutomaticEnv()
	return viper.GetInt(env)
}

func main() {
	fmt.Println(GetEnvInfo("MX_SHOPS"))
	//v := viper.New()
	////文件的路径如何设置
	////v.SetConfigName("config")
	////v.SetConfigType("yaml")
	////v.AddConfigPath("test01/viper_test/")
	//v.SetConfigFile("test01/viper_test/config.yaml")
	//if err := v.ReadInConfig(); err != nil {
	//	panic(err)
	//}
	//serverConfig := ServerConfig{}
	//if err := v.Unmarshal(&serverConfig); err != nil {
	//	panic(err)
	//}
	//fmt.Println(serverConfig)
	//fmt.Println(v.Get("name"))
	//
	//v.WatchConfig()
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//	_ = v.ReadInConfig()
	//	_ = v.Unmarshal(&serverConfig)
	//	fmt.Println(serverConfig)
	//})
	//time.Sleep(300 * time.Second)
}
