/**
* @Author vangogh
* @Description //TODO
* @File:  tencentuploader_test.go
* @Datetime 2022/5/9 17:07
**/
package mango_upload

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestTencentUploader_Upload(t *testing.T) {
	config := TencentConfig{
		SecretId:    "",
		SecretKey:   "",
		BucketURL:   "",
		ServiceURL:  "",
		UnzipPrefix: "unzip",
	}

	uploader, err := NewTencentUploader(config)
	if err != nil{
		t.Error(err)
	}

	// 读取压缩文件
	file, err := os.Open("test.zip")
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

	fmt.Println(uploader.Upload(fileData, "test.zip", UNZIP_FILE))
}