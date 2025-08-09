#### 查看go配置信息
	go env
#### 修改成国内的镜像
	go env -w GO111MODULE=on
	go env -w GOPROXY=https://goproxy.io,direct
### go安装import
	go get 