# gRPC demo

## 为什么有这个repo?
官方的demo竟然要clone整个grpc, 而且目录结构也很不舒服, 同时只有原生的部分, grpc拦截器和其他的几乎没有. 还是自己搞一个好了. 一是后面好复习, 现在刚好学习.

## 这里使用了什么
1. [gRPC](https://grpc.io/docs/quickstart/go.html)
2. [官方demo](https://github.com/grpc/grpc-go/tree/master/examples)
3. [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)
4. ~~[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)~~

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
将该repo的代码clone到`$GOPATH`中, 然后在proto目录下, 执行命令生成. 这里以helloworld为例
```sh
mkdir -p $GOPATH/src
git clone https://github.com/hxzhao527/grpcdemo.git $GOPATH/src/grpcdemo
cd $GOPATH/src/grpcdemo/proto
protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```
使用`go generate`也行, 自己看吧
### 4. 运行在app目录下的server和client
这里以helloworld为例
```sh
go generate $GOPATH/src/grpcdemo/pkg/service/helloworld
nohup go run $GOPATH/src/grpcdemo/app/server/main.go &
go run $GOPATH/src/grpcdemo/app/helloworld/client/main.go
```
如果要使用[Authentication](https://grpc.io/docs/guides/auth.html#go)
使用[lc-tlscert.go](https://raw.githubusercontent.com/driskell/log-courier/1.x/src/lc-tlscert/lc-tlscert.go)生成证书(可以解决[no-ip-sans](https://serverfault.com/questions/611120/failed-tls-handshake-does-not-contain-any-ip-sans))
将`selfsigned.crt`重命名为`public.pem`, `selfsigned.key`为`private.key`, 放到assets目录下. 在tools目录下有工具的拷贝, 可以使用
```sh
go generate $GOPATH/src/grpcdemo/pkg/service/helloworld
nohup go run $GOPATH/src/grpcdemo/app/server/main.go -ssl &
go run $GOPATH/src/grpcdemo/app/helloworld/client/main.go -ssl
```

## 有什么注意的吗?
1. ~~[openssl/1.1.1 bug](https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=898470)会导致`go generate`返回码不是0, 忽略就好.~~使用`lc-tlscert.go`就好了.
2. 为啥不用[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway), 鸡肋.不如使用`gin`写个web包一下还. 而且按照官方的指导竟然不能编译(可能是使用了dep的问题, 导致依赖默认拉不到最新).