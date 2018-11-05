# gRPC demo

## 为什么有这个repo?
官方的demo竟然要clone整个grpc, 而且目录结构也很不舒服, 同时只有简单的helloworld, grpc生态的其他都没有. 还是自己搞一个好了.

## 这里使用了什么
1. [gRPC](https://grpc.io/docs/quickstart/go.html)

## 怎么运行
### 1. 安装protoc compiler
可以去[protobuf](https://github.com/protocolbuffers/protobuf/releases)官方下zip, 也可以用系统的包管理工具, 比如arch下`sudo pacman -Syy protobuf`
### 2. 安装grpc的生成工具
protobuf只约定了数据格式, 自身工具并不能生成rpc调用相关的部分, 因此需要安装插件,
```sh
go get -u github.com/golang/protobuf/protoc-gen-go
export PATH=$PATH:$GOPATH/bin
```
注意最后将生成工具加到`PATH`中, 不然会报错
### 3. 生成代码
将该repo的代码clone到`$GOPATH`中, 然后在proto目录下, 执行命令生成.
```sh
mkdir -p $GOPATH/src
git clone https://github.com/hxzhao527/grpcdemo.git $GOPATH/src/grpcdemo
cd $GOPATH/src/grpcdemo/proto
protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```
### 4. 运行在app目录下的server和client
```sh
nohup go run $GOPATH/src/grpcdemo/app/server/main.go &
go run $GOPATH/src/grpcdemo/app/client/main.go
```