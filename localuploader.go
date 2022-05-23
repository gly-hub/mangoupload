/**
* @Author vangogh
* @Description 上传到本地
* @File:  uploadlocal
* @Datetime 2022/4/20 10:40
**/
package mango_upload

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/go-basic/ipv4"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type LocalConfig struct {
	Domain string // 访问地址域
	LocalPath string // 本地存储地址
	RoutePath string // 路由地址
}

type LocalUploader struct {
	domain string
	localPath string
	routePath string
}

func NewLocalUploader(config LocalConfig)(*LocalUploader){
	return &LocalUploader{
		domain:    config.Domain,
		localPath: config.LocalPath,
		routePath: config.RoutePath,
	}
}

func (lu *LocalUploader) Upload(fileData FileData, key string, fileType FileType)(string, error){
	var url string

	dirPath := fmt.Sprintf("%s/%s", lu.localPath, fileType)
	err := os.MkdirAll(dirPath, 0777)
	if err != nil {
		zap.L().Error("创建文件夹失败:", zap.Error(err))
		return "", err
	}

	filePath := fmt.Sprintf("%s/%s", lu.localPath, key)
	create, err := os.Create(filePath)
	if err != nil {
		zap.L().Error("创建文件失败:", zap.Error(err))
		return "", err
	}
	_, err = io.Copy(create, fileData.Body)
	if err != nil {
		zap.L().Error("构建文件失败:", zap.Error(err))
		return "", err
	}

	// 关闭
	_ = create.Close()

	ip := ipv4.LocalIP()
	if lu.domain != "" {
		ip = lu.domain
	}
	url = fmt.Sprintf("%s/%s/%s", ip, lu.routePath, key)

	// 当上传文件需要解压时，需要进行解压操作
	if fileType == UNZIP_FILE{
		// 解压
		err = lu.unzip(filePath, dirPath)
		if err != nil{
			return "", errors.New("unzip failed")
		}
		// 删除原压缩文件
		err  = os.Remove(filePath)
		if err != nil{
			fmt.Println(err)
			return "", errors.New("Failed to delete the zip file")
		}

		key = strings.ReplaceAll(key, path.Ext(key), "")
		url = fmt.Sprintf("%s/%s/%s/%s", ip, lu.routePath, fileType, key)
	}


	return url, nil
}

func (lu *LocalUploader) unzip(zipFile string,destDir string) error {
	zipReader ,err :=zip.OpenReader(zipFile)
	if err != nil{
		fmt.Println(err)
		zap.L().Error("读取压缩文件失败：", zap.Error(err))
		return err
	}
	defer zipReader.Close()
	for _,f := range zipReader.File{
		fPath := filepath.Join(destDir,f.Name)
		if f.FileInfo().IsDir(){
			os.MkdirAll(fPath,os.ModePerm)
		}else {
			if err = os.MkdirAll(filepath.Dir(fPath),os.ModePerm);err !=nil{
				return err
			}
			inFile,err := f.Open()  //这个是从压缩文件读取出来的
			if err !=nil{
				fmt.Println(err)
				zap.L().Error("读取文件失败：", zap.Error(err))
				return err
			}
			defer inFile.Close()

			outFile,err := os.OpenFile(fPath,os.O_CREATE|os.O_WRONLY|os.O_TRUNC,f.Mode()) //创建的新文件
			if err !=nil{
				fmt.Println(err)
				zap.L().Error("创建文件失败：", zap.Error(err))
				return err
			}
			defer outFile.Close()

			_,err = io.Copy(outFile,inFile)
			if err != nil{
				fmt.Println(err)
				zap.L().Error("构建文件失败：", zap.Error(err))
				return err
			}
		}
	}
	return err
}