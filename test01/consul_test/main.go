package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

func Register(address string, port int, name string, path string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.21.132:8500"

	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}
	// 生成注册对象
	errs := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			HTTP:                           "http://192.168.1.10:8021/health",
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if errs != nil {
		panic(errs)
	}
	return nil
}

func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.21.132:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	services, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for key, _ := range services {
		fmt.Println(key)
	}
}

func main() {
	//_ = Register("192.168.1.10", 8021, "user-web", "/health", []string{"mxshop", "boby"}, "user-web")
	AllServices()
}
