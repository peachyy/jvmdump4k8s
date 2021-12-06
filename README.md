# jvm dump 在k8s中导出OSS工具

https://www.cnblogs.com/peachyy/p/15539217.html


#### 现状

加参数 -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=logs/test.dump 可以实现在jvm发生内存错误后 会生成dump文件 方便开发人员分析异常原因。

当运行在k8s中，如果进程发生错误 导出dump文件后 ，k8s会重启dokcer容器，上一次崩溃生成的dump文件就没有了。如果应用并没有完全崩溃 此时极其不稳定 最好也能通知到技术人员来处理。这样不方便我们排查原因 所有写了一个小工具。大概原理如下

1、 -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=logs/test.dump 当发生内存错误的时候 导出堆文件
2、 -XX:OnOutOfMemoryError=./dumpError.sh 当发生内存溢出的时候，让JVM调用一个shell脚本 这个shell脚本可以做一些资源整理操作 比如kill掉当前进程并重启

依赖上面2点jvm特性 就能做到把dump文件收集起来 是通知技术人员也好(比如发送订单、短信报警等)、然后再把dump文件上传到OSS 或者其他的文件存储中。 需要值得注意的是-XX:OnOutOfMemoryError=xx.sh 执行的脚本不能传脚本参数，所以尽可能把参数都封装在另一个脚本中。


了解更多 https://www.cnblogs.com/peachyy/p/15539217.html


#### 打包构建

window

```
  go build -ldflags "-w -s"
```

linux
```
  set GOOS=linux
  go build -ldflags "-w -s"
```

#### 使用示例
```
wget https://github.com/peachyy/jvmdump4k8s/releases/download/v1.11.2/jvmdump4k8s-linux.zip
unzip jvmdump4k8s-linux.zip 
chmod +x jvmdump4k8s

 vim jvmdump4k8s.ini 

```
输入相关配置并保存
storage.type 可选值 阿里云OSS=alioss 华为云OBS=huaweiobs

```
storage.type=alioss

#阿里云OSS配置
alioss.endpoint=请输入ossendpoint
alioss.accessKey=请输入
alioss.accessSecret=请输入
alioss.bucketName=请输入
#oss目录
alioss.folder=jvmdump

#华为云OBS配置
huaweiobs.endpoint=请输入obs  endpoint
huaweiobs.ak=EOQALVN6RH4R0K7PRKYS
huaweiobs.sk=4caa7n0QyLGzzbAyXcxwYQSP2XmDvX6HTYGOnP49
huaweiobs.bucketName=changan-obs-app2
huaweiobs.folder=jvmdump

#通知消息是可选的 如果不设置就不推送消息
#钉钉群机器人token
notify.dingtalkToken=
#微信群webhook key
notify.wxKey=
#可以在配置文件中指定dump文件地址 也可以在命令行用-f指定 会覆盖此配置
dump.filepath=
```
测试一下配置与上传 

-f 指定Dump文件路径 需要大于5M 这里使用的是分片上传。一般dump文件都会大于5m

-pod 指定pod的名称 用于消息推送的备注是哪个Pod触发的OOM

```
 ./jvmdump4k8s  -f dump文件路径 
```

如果能成功上传就表示配置这些没有问题。接下来可以集成到jvm中去。

使用以下命令在JVM发生OOM的时候 会自动生成dump文件在logs/manager.dump中，需要特殊说明一下 `OnOutOfMemoryError` 
要用一个脚本包装一下 因为不支持传参数，参数都放在shell中传给jvmdump4k8s工具。

```
 -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=logs/manager.dump  -XX:OnOutOfMemoryError=./dumpError.sh
```

dumpError.sh 脚本内容大概如下，pod名称就用`$HOSTNAME`获取 

```
#!/bin/bash
./jvmdump4k8s -f logs/manager.dump -pod $HOSTNAME
```
