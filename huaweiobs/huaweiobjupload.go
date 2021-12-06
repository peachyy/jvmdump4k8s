package huaweiobs

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"jvmdump4k8s/config"
	"jvmdump4k8s/util"
	"os"
	"path"
	"path/filepath"
)

func UploadToHwObs(file string) string {
	obsClient, err := obs.New(config.GlobalConfig.HwAk, config.GlobalConfig.HwSk, config.GlobalConfig.HwEndpoint)
	if err != nil {
		panic(err)
	}
	bucketName := config.GlobalConfig.HwBucketName
	var filename = filepath.Base(file) //获取文件名
	var ext = path.Ext(file)           //获取扩展名
	ossFileName := config.GlobalConfig.HwFolder + "/" + filename + util.FormartdateNow() + ext
	// Claim a upload id firstly
	input := &obs.InitiateMultipartUploadInput{}
	input.Bucket = bucketName
	input.Key = ossFileName
	output, err := obsClient.InitiateMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	uploadId := output.UploadId

	fmt.Printf("Claiming a new upload id %s\n", uploadId)
	fmt.Println()

	// Calculate how many blocks to be divided
	// 5MB
	var partSize int64 = 5 * 1024 * 1024

	stat, err := os.Stat(file)
	if err != nil {
		panic(err)
	}
	fileSize := stat.Size()

	partCount := int(fileSize / partSize)

	if fileSize%partSize != 0 {
		partCount++
	}
	fmt.Printf("Total parts count %d\n", partCount)
	fmt.Println()

	//  Upload parts
	fmt.Println("Begin to upload parts to OBS")

	partChan := make(chan obs.Part, 5)

	for i := 0; i < partCount; i++ {
		partNumber := i + 1
		offset := int64(i) * partSize
		currPartSize := partSize
		if i+1 == partCount {
			currPartSize = fileSize - offset
		}
		go func() {
			uploadPartInput := &obs.UploadPartInput{}
			uploadPartInput.Bucket = bucketName
			uploadPartInput.Key = ossFileName
			uploadPartInput.UploadId = uploadId
			uploadPartInput.SourceFile = file
			uploadPartInput.PartNumber = partNumber
			uploadPartInput.Offset = offset
			uploadPartInput.PartSize = currPartSize
			uploadPartInputOutput, errMsg := obsClient.UploadPart(uploadPartInput)
			if errMsg == nil {
				fmt.Printf("part %d finished\n", partNumber)
				partChan <- obs.Part{ETag: uploadPartInputOutput.ETag, PartNumber: uploadPartInputOutput.PartNumber}
			} else {
				panic(errMsg)
			}
		}()
	}

	parts := make([]obs.Part, 0, partCount)

	for {
		part, ok := <-partChan
		if !ok {
			break
		}
		parts = append(parts, part)
		if len(parts) == partCount {
			close(partChan)
		}
	}

	fmt.Println()
	fmt.Println("Completing to upload multiparts")
	completeMultipartUploadInput := &obs.CompleteMultipartUploadInput{}
	completeMultipartUploadInput.Bucket = bucketName
	completeMultipartUploadInput.Key = ossFileName
	completeMultipartUploadInput.UploadId = uploadId
	completeMultipartUploadInput.Parts = parts

	//sample.doCompleteMultipartUpload(completeMultipartUploadInput)
	cmpleteFile, err2 := obsClient.CompleteMultipartUpload(completeMultipartUploadInput)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("Complete multiparts finished")
	setAclRead(*obsClient, ossFileName)
	url := bucketName + "." + config.GlobalConfig.HwEndpoint + "/" + cmpleteFile.Key
	fmt.Println("url ", url)
	return url
}

//func getSignUrl(obsClient obs.ObsClient,filekey string)  {
//	input := &obs.CreateSignedUrlInput{}
//	input.Method = obs.HttpMethodPut
//	input.Bucket = bucketName
//	input.Key = filekey
//	input.SubResource = obs.SubResourceAcl
//	//input.Expires = 3600
//	input.Headers = map[string]string{obs.HEADER_ACL_AMZ: string(obs.AclPublicRead)}
//	output, err := obsClient.CreateSignedUrl(input)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("%s using temporary signature url:\n", "SetObjectAcl")
//	fmt.Println(output.SignedUrl)
//}
//设置为公共读
func setAclRead(obsClient obs.ObsClient, filekey string) {
	input := &obs.SetObjectAclInput{}
	input.Bucket = config.GlobalConfig.HwBucketName
	input.Key = filekey
	input.ACL = obs.AclType("public-read")
	output, err := obsClient.SetObjectAcl(input)
	if err == nil {
		fmt.Printf("setAclRead RequestId success:%s\n", output.RequestId)
	} else {
		fmt.Printf("setAclRead RequestId ERROR:%s\n", err)
	}
}
