/**
* @Author vangogh
* @Description 腾讯cos上传器
* @File:  tencentuploader
* @Datetime 2022/5/9 16:23
**/
package mango_upload

import (
	"context"
	"errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type TencentConfig struct {
	SecretId string // secretId
	SecretKey string // secretKey
	BucketURL string // 桶URL
	ServiceURL string // 服务URL
	UnzipPrefix string // 解压文件前缀
}

type TencentUploader struct {
	cosClient *cos.Client
	unzipPrefix string
	bucketUrl string
}

func NewTencentUploader(config TencentConfig)(uploader *TencentUploader, err error){
	var (
		u *url.URL
		su *url.URL
		cosClient *cos.Client
		ok bool
	)
	u, err = url.Parse(config.BucketURL)
	if err != nil{
		return
	}

	su, err = url.Parse(config.ServiceURL)
	b := &cos.BaseURL{BucketURL:u, ServiceURL:su}
	cosClient = cos.NewClient(b, &http.Client{
		Transport:     &cos.AuthorizationTransport{
			SecretID:     config.SecretId,
			SecretKey:    config.SecretKey,
		},
	})

	ok, err = cosClient.Bucket.IsExist(context.Background())
	if err != nil{
		return
	}
	if !ok {
		err = errors.New("bucket not exists")
		return
	}

	uploader = &TencentUploader{
		cosClient: cosClient,
		unzipPrefix: config.UnzipPrefix,
		bucketUrl: config.BucketURL,
	}
	return
}

func (tcu *TencentUploader) Upload(fileData FileData, key string, fileType FileType) (string, error){
	newKey := key
	if fileType == UNZIP_FILE {
		newKey = tcu.unzipPrefix + ":" + key

		key = strings.ReplaceAll(key, path.Ext(newKey), "")
		_, err := tcu.cosClient.Object.Delete(context.Background(), key)
		if err != nil{
			return "", err
		}
	}

	_, err := tcu.cosClient.Object.Put(context.Background(), newKey, fileData.Body, nil)
	if err != nil{
		return "", err
	}

	urlT := tcu.bucketUrl + "/" + key
	return urlT, nil
}