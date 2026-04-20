# 安装工具
    
    https://github.com/protocolbuffers/protobuf/releases

    go get google.golang.org/protobuf/protoc-gen-go


# 生成proto文件
    
    //老版本
    protoc -I . helloworld.proto --go_out=plugins=grpc:.
    //新版本
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helloworld.proto
