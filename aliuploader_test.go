/**
* @Author vangogh
* @Description //TODO
* @File:  aliuploader_test.go
* @Datetime 2022/4/21 14:28
**/
package mango_upload

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestALiUploader_Upload(t *testing.T) {
	config := ALiConfig{
		AccessKey:  "",
		AccessId:  "",
		Endpoint:   "",
		BucketName: "",
		Domain: "",
		UnzipPrefix: "",
		UnzipDir: "",
	}

	uploader, err := NewALiUploader(config)
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