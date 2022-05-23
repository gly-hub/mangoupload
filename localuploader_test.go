/**
* @Author vangogh
* @Description //TODO
* @File:  lcoaluploader_test.go
* @Datetime 2022/4/20 15:09
**/
package mango_upload

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestLocalUploader_Upload(t *testing.T) {
	// 初始化本地上传器
	localUploader := LocalUploader{
		domain:    "https://test.com",
		localPath: "C:/Users/vangogh/Desktop/upload/local",
		routePath: "v1/local",
	}

	// 读取压缩文件
	file, err := os.Open("C:/Users/vangogh/Desktop/mango.zip")
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

	// 上传测试
	fmt.Println(localUploader.Upload(fileData, "mango.zip", UNZIP_FILE))
}
