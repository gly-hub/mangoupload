/**
* @Author vangogh
* @Description 七牛上传器
* @File:  qiniuuploader
* @Datetime 2022/4/21 15:26
**/
package mango_upload

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiNiuConfig struct {
	Domain string
	AccessKey string
	SecretKey string
	BucketName string
}

type QiNiuUploader struct {
	qBoxClient *qbox.Mac
	domain string
	bucketName string
}

func NewQiNiuUploader(config QiNiuConfig)(uploader *QiNiuUploader, err error){
	var qBoxClient *qbox.Mac
	qBoxClient = qbox.NewMac(config.AccessKey, config.SecretKey)

	bm := storage.NewBucketManager(qBoxClient, &storage.Config{})
	_, err = bm.GetBucketInfo(config.BucketName)
	if err != nil{
		return
	}

	uploader = &QiNiuUploader{
		qBoxClient: qBoxClient,
		domain:     config.Domain,
		bucketName: config.BucketName,
	}

	return
}

func (qnu *QiNiuUploader) Upload(fileData FileData, key string, fileType FileType)(string, error){
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", qnu.domain, key),
	}
	upToken := putPolicy.UploadToken(qnu.qBoxClient)

	cfg := storage.Config{
		Zone:          &storage.Zone{},
		UseHTTPS:      true,
		UseCdnDomains: true,
	}

	// 上传
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.Put(context.Background(), &ret, upToken, key, fileData.Body, fileData.Size, &storage.PutExtra{})
	if err != nil{
		return "", err
	}

	url := qnu.domain + "/" + ret.Key

	return url, nil
}