# gRPC demo

## What's this?
As the title says, this is a demo for [gRPC](https://grpc.io/docs/quickstart/go.html)<br>
如题, 一个grpc的demo.

## why not official example?
The [official example](https://github.com/grpc/grpc-go/tree/master/examples) contains:
1. [Unary RPC](https://grpc.io/docs/guides/concepts.html#unary-rpc) in [HelloWorld](https://github.com/grpc/grpc-go/tree/master/examples/helloworld).
2. [Stream request](https://grpc.io/docs/guides/concepts.html#server-streaming-rpc) in [route_guide](https://github.com/grpc/grpc-go/tree/master/examples/route_guide)
3. [Metadata and Interceptor](https://grpc.io/docs/guides/concepts.html#metadata) in [oauth](https://github.com/grpc/grpc-go/tree/master/examples/oauth)
4. [Error model](https://grpc.io/docs/guides/error.html) in [rpc_errors](https://github.com/grpc/grpc-go/tree/master/examples/rpc_errors)

It basically covers all the basic concepts in gRPC.But that is **not enough**. It lacks of the example about integration with other components, e.g [Load Balancing](https://github.com/grpc/grpc/blob/master/doc/load-balancing.md), [Name Resolution](https://github.com/grpc/grpc/blob/master/doc/naming.md).And its file structure doesn't have any helpful instuction. So I restruct it according to the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) and add some codes to extend.

And this repo is my learning notes. It will be helpful when I use gRPC at work later.

官方的demo有点简陋, 虽然涵盖了grpc基本的点, 但是离实际使用有点远或者说缺少太多生产中关注的东西. 比如LB, 比如名字解析.这些东西在实际项目中很重要, 但是相关的资料搜起来有点费劲, 而且搜到的多是旧版本,官方已不推荐的. 还有就是官方的文件组织有点蛋疼, 实际项目完全无法借鉴. 于是我参考[golang-standards/project-layout](https://github.com/golang-standards/project-layout)重新组织了一下, 然后又添了一点周边扩展的代码在里面.

当然, 这是我刚开始学gRPC, 开这个repo也是为了做做笔记, 后面真在项目里用的时候也好复习.

## What this repo contains?
1. [offical example](https://github.com/grpc/grpc-go/tree/master/examples) after restructed
2. [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)
3. ~~[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)~~ :( I will explain later.
4. integration with [consul](https://www.consul.io)
5. [Authentication](https://grpc.io/docs/guides/auth.html#go)
6. [Health Checking](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)

目前为止, 项目里包含里以下内容:
1. 重新组织后的官方示例
2. [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware) 项目
3. 与[consul](https://www.consul.io)的e整合
4. **没有[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)**, 原因后面说
5. [证书](https://grpc.io/docs/guides/auth.html#go) 或叫https, tls
6. [健康检查](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)

## Dependency
1. [dep](https://github.com/golang/dep)
2. [protocolbuffers compiler](https://github.com/protocolbuffers/protobuf), you can download from [releases](https://github.com/protocolbuffers/protobuf/releases) or package-manager you like
3. [the grpc plugin for protoc](https://github.com/golang/protobuf/tree/master/protoc-gen-go)

Make sure the executables  `dep`, `protoc`, `protoc-gen-go` in your `${PATH}`.
If you install `dep` and `protoc-gen-go` by go toolchain, you may need to add these by `export PATH=$PATH:$GOPATH/bin`.

项目是使用[dep](https://github.com/golang/dep)做依赖管理的, 因此需要机器上只有可运行的dep. 同时为了能根据`proto`文件正常生成代码, 还需要`protoc`编译器及其`grpc`插件.
自己按照官方install安装就好了, 不过如果系统允许, 包管理器也不错. 比如arch下 可以使用`sudo pacman -Syy protobuf`安装`protocolbuffers compiler`

最后, 确保`dep`, `protoc`, `protoc-gen-go`在你的环境变量里, 因为后面会用的到.

## How to run?
Take `HelloWorld` as the example:
```sh
export GOPATH=/home/hxzhao/goproject # <= $GOPATH, wherever you like
mkdir -p $GOPATH/src
git clone https://github.com/hxzhao527/grpcdemo.git $GOPATH/src/grpcdemo # clone code to local
go generate $GOPATH/src/grpcdemo/internal/service/helloworld
nohup go run $GOPATH/src/grpcdemo/app/server/main.go &
go run $GOPATH/src/grpcdemo/app/client/helloworld/main.go
```
The `server` will run in the background. After the client return, you can just terminate it.

If you want to secure your server with [Authentication](https://grpc.io/docs/guides/auth.html#go).
```sh
go generate $GOPATH/src/grpcdemo/app/server
nohup go run $GOPATH/src/grpcdemo/app/server/main.go -ssl &
go run $GOPATH/src/grpcdemo/app/client/helloworld/main.go -ssl
```
*Because of [openssl/1.1.1 bug](https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=898470), `go generate` maybe print some error messages, just ignore. Or, use the `cert-maker` in [tools](./tools)*

About the other startup parameter, code is here, just read and find those by youself.

以`HelloWorld`为例, 可以直接运行上面的脚本, 如果想用证书加密client与server的通讯, 可以先生成一下证书, 然后启动时加`-ssl`参数. 至于其他的启动参数, 翻翻代码就知道了.

由于openssl的bug, 生成证书时可能会输出一下报错信息, 忽略就好.

## How to contribute?
I'm gland to see your PR. If you have any question about this repo or gRPC, just commit an issue.

这毕竟是我边学习边写的, 难免有写的烂的地方. 如果你有什么更好的实现或者觉得这代码屎一样, 兄弟亮出你的代码,咱PR见. 如果你在使用这个demo或者使用grpc时遇到任何问题, 欢迎提issue, 咱可以一起讨论一起研究.

## Appendix
1. why no [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway). It failed me, so sad. And I prefer [envoy](https://www.envoyproxy.io/).
