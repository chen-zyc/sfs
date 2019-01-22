# sfs

简单的静态文件服务器

在启动脚本的目录下启动一个简单的文件服务，可以访问当前目录下的任意文件。

比如在 `~/tmp/sfs` 下启动：

```shell
~/tmp/sfs > sfs
Run file server(dir=/Users/yuchenzhang/tmp/sfs) on 127.0.0.1:9000
```

在 `~/tmp/sfs` 下有一个名为 hello 的文件，访问它可以通过：

```shell
curl http://127.0.0.1:9000/hello
```

指定服务地址： `sfs -addr=127.0.0.1:9000`



## 打印请求信息

dump 请求头： `sfs -dump=1`

```shell
==================== Incoming Request ====================
GET /hello HTTP/1.1
Host: localhost:9000
Accept: */*
User-Agent: curl/7.54.0
==================== Dump Finished ====================
```

dump 请求头和请求体： `sfs -dump=2`

```shell
==================== Incoming Request ====================
POST /hello HTTP/1.1
Host: localhost:9000
Accept: */*
Content-Length: 12
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/7.54.0

request body
==================== Dump Finished ====================
```



## 处理请求前阻塞一段时间

指定在处理请求之前阻塞多长时间： 

```bash
sfs -block-header="Block"
```

block-header 指定由哪个请求头指定阻塞时间。 比如：

```bash
curl "http://localhost:9000/hello" -H "Block: 3s"
```

