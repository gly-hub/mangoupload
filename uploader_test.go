/**
* @Author vangogh
* @Description //TODO
* @File:  uploader_test.go
* @Datetime 2022/4/20 20:13
**/
package mango_upload

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestUploader_Upload(t *testing.T) {
	uploader := NewUploader()

	uploader.SetKeyModel(true)

	localUploader := NewLocalUploader(LocalConfig{
		Domain:    "",
		LocalPath: "",
		RoutePath: "",
	})

	uploader.RegisterUploader(localUploader)

	//
	//uploader.Upload(testFile, IMAGE)
}

func TestUploader_Md5(t *testing.T) {
	uploader := NewUploader()
	fmt.Println(hex.EncodeToString((uploader.Md5([]byte("123")))))
}