package config

import (
	"github.com/go-ini/ini"
)

type Config struct {
	Type         string `ini:"storage.type"`
	Endpoint     string `ini:"alioss.endpoint"`
	AccessKey    string `ini:"alioss.accessKey"`
	AccessSecret string `ini:"alioss.accessSecret"`
	BucketName   string `ini:"alioss.bucketName"`
	Folder       string `ini:"alioss.folder"`

	NotifyDingToken string `ini:"notify.dingtalkToken"`
	NotifyWxToken   string `ini:"notify.wxKey"`
	PodName         string `ini:"podName"`
	DumpFile        string `ini:"dump.filepath"`
}

var inifile = "jvmdump4k8s.ini"
var GlobalConfig Config = Config{}

func init() {
	//加载INI文件
	cfg, err := ini.Load(inifile)
	err = cfg.Section("").MapTo(&GlobalConfig)
	if err != nil {
		panic(err)
	}
}
