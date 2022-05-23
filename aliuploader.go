/**
* @Author vangogh
* @Description 阿里云oss上传器
* @File:  aliuploader
* @Datetime 2022/4/21 14:08
**/
package mango_upload

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
	"path"
	"strings"
)

type ALiConfig struct {
	AccessId string
	AccessKey string
	Endpoint string
	BucketName string
	Domain string
	UnzipPrefix string // 解压前缀
	UnzipDir string // 解压到目录
}


type ALiUploader struct {
	ossClient *oss.Client
	bucketName string
	domain string
	unzipPrefix string
	unzipDir string
}

func NewALiUploader(config ALiConfig)(uploader *ALiUploader, err error){
	var ossClient *oss.Client
	ossClient, err = oss.New(config.Endpoint, config.AccessId, config.AccessKey)
	if err != nil{
		return
	}

	// 校验桶
	_, err = ossClient.Bucket(config.BucketName)
	if err != nil{
		return
	}

	uploader = &ALiUploader{
		ossClient: ossClient,
		bucketName: config.BucketName,
		domain: config.Domain,
		unzipPrefix: config.UnzipPrefix,
		unzipDir: config.UnzipDir,
	}
	return
}

func (alu *ALiUploader) Upload(fileData FileData, key string, fileType FileType) (string, error){
	url := alu.domain + "/" + key
	newKey := key
	bucket, err := alu.ossClient.Bucket(alu.bucketName)
	if err != nil {
		return "", err
	}

	if fileType == UNZIP_FILE {
		newKey = alu.unzipPrefix + key
		key = strings.ReplaceAll(key, path.Ext(newKey), "")
		lor, err := bucket.ListObjects()
		if err != nil {
			zap.L().Error("获取存储桶对象错误：", zap.Error(err))
		}

		objects := make([]string, 0)
		prefix := fmt.Sprintf("%s/%s/", alu.unzipDir, key)
		for _, content := range lor.Objects{
			if len(content.Key) > len(prefix) && prefix == content.Key[:len(prefix)]{
				objects = append(objects, content.Key)
			}
		}

		if len(objects) > 0 {
			_, err = bucket.DeleteObjects(objects, oss.DeleteObjectsQuiet(false))
			if err != nil {
				zap.L().Error("删除存储桶对象失败:", zap.Error(err))
			}
		}

		url = alu.domain + "/" +  alu.unzipDir + "/" + key
	}

	// 上传
	err = bucket.PutObject(newKey, fileData.Body)
	if err != nil {
		zap.L().Error("上传对象失败:", zap.Error(err))
		return "", err
	}

	return url, nil
}