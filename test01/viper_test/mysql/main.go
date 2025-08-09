package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Name      string      `mapstructure:"name"`
	MysqlInfo MysqlConfig `mapstructure:"mysql"`
}

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func main() {
	v := viper.New()
	//文件的路径如何设置
	//v.SetConfigName("config")
	//v.SetConfigType("yaml")
	//v.AddConfigPath("test01/viper_test/")
	v.SetConfigFile("test01/viper_test/mysql/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig.MysqlInfo.Port)
	//fmt.Println(v.Get("name"))
}
