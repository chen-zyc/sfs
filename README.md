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