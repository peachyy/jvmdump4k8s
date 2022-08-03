package main

import (
	"flag"
	"fmt"
	"jvmdump4k8s/alioss"
	"jvmdump4k8s/config"
	"jvmdump4k8s/huaweiobs"
	"jvmdump4k8s/notify"
	"jvmdump4k8s/qiniu"
	"jvmdump4k8s/util"
)

var (
	enabledd bool   //是否启用通知
	enablewx bool   //是否启用微信通知
	dumpFile string //dump的文件
	pod      string //podname

)

func init() {
	//fmt.Println("init....")

}

func main() {
	parseCli()
	fmt.Println("start invoke dump...")
	flag.Parse()
	if "" != pod {
		config.GlobalConfig.PodName = pod
	}
	if "" != dumpFile {
		config.GlobalConfig.DumpFile = dumpFile
	}
	fmt.Printf("dumpFile %s \n", dumpFile)
	//dump文件是否存在
	exist, err := util.FileExists(dumpFile)
	if err != nil {
		fmt.Printf("验证文件是否存在发生错误![%v]\n", err)
		return
	}
	if exist {
		var url = uploadStorage(dumpFile)
		fmt.Printf("OSS上传完成 %s\n", url)
		notifyFunc(url)
	} else {
		fmt.Printf("dump文件不存在 %s\n", dumpFile)
	}
}

//解析命令行参数
func parseCli() {

	flag.StringVar(&dumpFile, "f", "", "-f xx.dump")
	flag.StringVar(&pod, "pod", "", "pod")
	if "" != config.GlobalConfig.NotifyDingToken {
		enabledd = true
	}
	if "" != config.GlobalConfig.NotifyWxToken {
		enablewx = true
	}
}

//Storage
func uploadStorage(file string) string {
	type_ := config.GlobalConfig.Type
	switch type_ {
	case "alioss":
		return alioss.Upload(file)
	case "huaweiobs":
		return huaweiobs.UploadToHwObs(file)
	case "qiniu":
		return qiniu.Upload(file)
	default:
		panic(fmt.Sprintf("不支持文件存储类型%s", type_))
	}
}

//发送IM 工具
func notifyFunc(fileurl string) {
	if enabledd {
		notify.SendDingtalk(fileurl)
	}
	if enablewx {
		notify.SendWechat(fileurl)
	}

}
