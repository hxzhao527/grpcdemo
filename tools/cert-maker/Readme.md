# 证书制作工具
如果想要签ip的证书, 使用`openssl`, 有点麻烦. 可以参考[Using SSL with an IP address instead of DNS](https://bowerstudios.com/node/1007).
当然也可以使用这个目录的工具, `lc-tlscert`. 这是大佬[Jason Woods](https://github.com/driskell)写的一个工具, 原文件在[lc-tlscert.go](https://github.com/driskell/log-courier/blob/master/lc-tlscert/lc-tlscert.go), 这里拷贝一份留存.

# 为什么会涉及到签IP?
服务IP直连, 这时候用ssl就要签IP, 不然可能会出现[no-ip-sans](https://serverfault.com/questions/611120/failed-tls-handshake-does-not-contain-any-ip-sans).
