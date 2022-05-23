/**
* @Author vangogh
* @Description 华为上传器
* @File:  huaweiuploader
* @Datetime 2022/4/20 15:42
**/
package mango_upload

import (
	"errors"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"go.uber.org/zap"
	"path"
	"strings"
)

type HuaWeiConfig struct {
	AccessKey string // accessKey
	SecretKey string // secretKey
	Endpoint string // endPoint
	BucketName string // 桶名
	Domain string // 访问域
	UnzipPrefix string // 解压前缀
}

type HuaWeiUploader struct {
	obsClient *obs.ObsClient
	bucketName string
	domain string
	unzipPrefix string
}

func NewHuaWeiUploader(config HuaWeiConfig)(uploader *HuaWeiUploader, err error){
	var obsClient *obs.ObsClient
	obsClient, err = obs.New(config.AccessKey, config.SecretKey, config.Endpoint)
	if err != nil{
		return
	}

	// 校验桶
	_, err = obsClient.HeadBucket(config.BucketName)
	if err != nil{
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.StatusCode == 404 {
				err = errors.New("Bucket does not exists")
				return
			}
		}
	}

	uploader = &HuaWeiUploader{
		obsClient:obsClient,
		bucketName: config.BucketName,
		domain: config.Domain,
		unzipPrefix: config.UnzipPrefix,
	}
	return
}

func (hwu *HuaWeiUploader) Upload(fileData FileData, key string, fileType FileType) (string, error) {
	newKey := key
	if fileType == UNZIP_FILE {
		newKey = hwu.unzipPrefix + key

		input := &obs.ListObjectsInput{}
		input.Bucket = hwu.bucketName
		output, err := hwu.obsClient.ListObjects(input)
		if err != nil {
			zap.L().Error("获取存储桶对象错误：", zap.Error(err))
		}
		key = strings.ReplaceAll(key, path.Ext(newKey), "")

		objects := make([]obs.ObjectToDelete, 0, len(output.Contents))
		for _, content := range output.Contents {
			if len(content.Key) > len(key) && key == content.Key[:len(key)]{
				objects = append(objects, obs.ObjectToDelete{Key: content.Key})
			}
		}
		if len(objects) > 0{
			deleteObjectsInput := &obs.DeleteObjectsInput{}
			deleteObjectsInput.Bucket = hwu.bucketName
			deleteObjectsInput.Objects = objects[:]
			_, err = hwu.obsClient.DeleteObjects(deleteObjectsInput)
			if err != nil{
				zap.L().Error("刪除存储桶对象错误：", zap.Error(err))
			}
		}
	}
	input := &obs.PutObjectInput{}
	input.Bucket = hwu.bucketName
	input.Key = newKey
	input.Body = fileData.Body

	// 上传
	output, err := hwu.obsClient.PutObject(input)
	if err == nil {
		fmt.Printf("upload RequestId:%s\n", output.RequestId)
	} else if obsError, ok := err.(obs.ObsError); ok {
		err = errors.New(obsError.Message)
		zap.L().Error("上传对象失败:", zap.Error(err))
		return "", err
	}

	url := hwu.domain + "/" + key
	return url, nil
}