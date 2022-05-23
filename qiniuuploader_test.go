/**
* @Author vangogh
* @Description //TODO
* @File:  qiniuuploader_test.go
* @Datetime 2022/5/9 16:21
**/
package mango_upload

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestQiNiuUploader_Upload(t *testing.T) {
	config := QiNiuConfig{
		AccessKey:  "",
		SecretKey:  "",
		BucketName: "",
		Domain: "",
	}

	uploader, err := NewQiNiuUploader(config)
	if err != nil{
		t.Error(err)
	}

	// 读取压缩文件
	file, err := os.Open("")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileSize := fileinfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	var fileData = FileData{
		Body: bytes.NewBuffer(buffer),
		Size: fileSize,
	}

	fmt.Println(uploader.Upload(fileData, "", UNZIP_FILE))
}