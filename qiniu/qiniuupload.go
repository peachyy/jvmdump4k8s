package qiniu

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"jvmdump4k8s/config"
	"jvmdump4k8s/util"
	"os"
	"path"
	"path/filepath"
)

func Upload(file string) string {

	var bucketName = config.GlobalConfig.QiniuBucket
	var accessKey = config.GlobalConfig.QiniuAccessKey
	var accessSecret = config.GlobalConfig.QiniuSecretKey
	var folder = config.GlobalConfig.QiniuFolder
	var apiHost = config.GlobalConfig.QiniuApiHost
	fmt.Printf("开始上传七牛云OSS accessKey=%s bucketName=%s \n", accessKey, bucketName)
	putPolicy := storage.PutPolicy{
		Scope: bucketName,
	}
	mac := qbox.NewMac(accessKey, accessSecret)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	//cfg.Zone = &storage.ZoneHuadong
	//cfg.Zone=&storage.ZoneHuanan
	cfg.ApiHost = apiHost
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = true
	resumeUploader := storage.NewResumeUploaderV2(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.RputV2Extra{}
	var filename = filepath.Base(file) //获取文件名
	var ext = path.Ext(file)           //获取扩展名
	var objectName = folder + "/" + filename + util.FormartdateNow() + ext
	err := resumeUploader.PutFile(context.Background(), &ret, upToken, objectName, file, &putExtra)
	if err != nil {
		fmt.Println("七牛文件上传发生错误", err)
		os.Exit(-1)
		return ""
	}
	url := apiHost + "/" + objectName
	fmt.Printf("上传成功 %s\n", url)
	fmt.Println(ret.Key, ret.Hash)

	return url
}
