package alioss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"jvmdump4k8s/config"
	"jvmdump4k8s/util"
	"os"
	"path"
	"path/filepath"
)

//https://help.aliyun.com/document_detail/88604.html
//阿里云存储上传
func Upload(file string) string {
	fmt.Printf("开始上传阿里云OSS %s \n", file)
	var bucketName = config.GlobalConfig.BucketName
	var endpoint = config.GlobalConfig.Endpoint //外网：oss-cn-hangzhou.aliyuncs.com 内网：oss-cn-hangzhou-internal.aliyuncs.com
	var accessKey = config.GlobalConfig.AccessKey
	var accessSecret = config.GlobalConfig.AccessSecret
	var folder = config.GlobalConfig.Folder
	//if env == "test" {
	//	bucketName = "请设置测试buckname"
	//
	//	accessKey = "请设置测试 accessKey"
	//	accessSecret = "请设置测试 accessSecret"
	//} else {
	//	bucketName = "请设置正式 bucketName"
	//	accessKey = "请设置正式 accessKey"
	//	accessSecret = "请设置正式 accessSecret"
	//}

	client, err := oss.New(endpoint, accessKey, accessSecret) //建议oss内网地址[需要修改]
	if err != nil {
		fmt.Println("创建ali OSS连接失败:", err)
		os.Exit(-1)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("获取bucket失败:", err)
		os.Exit(-1)
	}
	// 将本地文件分片，且分片数量指定为3。
	chunks, err := oss.SplitFileByPartNum(file, 3)
	fd, err := os.Open(file)
	defer fd.Close()
	var filename = filepath.Base(file) //获取文件名
	var ext = path.Ext(file)           //获取扩展名
	// 设置存储类型为标准存储。
	storageType := oss.ObjectStorageClass(oss.StorageStandard)

	var objectName = folder + "/" + filename + util.FormartdateNow() + ext
	// 步骤1：初始化一个分片上传事件，并指定存储类型为标准存储。
	fmt.Printf("objectName=%s\n", objectName)
	imur, err := bucket.InitiateMultipartUpload(objectName, storageType)
	// 步骤2：上传分片。
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		fd.Seek(chunk.Offset, os.SEEK_SET)
		// 调用UploadPart方法上传每个分片。
		part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		parts = append(parts, part)
	}
	// 指定Object的读写权限为公共读，默认为继承Bucket的读写权限。
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	// 步骤3：完成分片上传，指定文件读写权限为公共读。
	cmur, err := bucket.CompleteMultipartUpload(imur, parts, objectAcl)
	if err != nil {
		fmt.Println("分片上传失败:", err)
		os.Exit(-1)
	}
	fmt.Println("完成分片上传:", cmur)
	return cmur.Location
}
