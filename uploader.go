/**
* @Author vangogh
* @Description 上传器
* @File:  uploader
* @Datetime 2022/4/20 10:42
**/
package mango_upload

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"mime/multipart"
	"path"
	"strconv"
)

type UploadImageResp struct {
	Url string `json:"url"`
}

type FileData struct {
	Body io.Reader
	Size int64
}

type IUploader interface {
	Upload(data FileData, key string, fileType FileType) (string, error)
}

// 上传器
type Uploader struct {
	md bool
	uploader IUploader // 本地上传器
	fileSizeConfig string
	fileExtConfig []string
}

func NewUploader() (*Uploader) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}
	zap.ReplaceGlobals(logger)
	return &Uploader{
		fileSizeConfig: FileSizeConfig,
		fileExtConfig: FileExtConfig,
	}
}

func (u *Uploader) SetUploadLimit(fileSizeConfig string, fileExtConfig []string){
	u.fileExtConfig = fileExtConfig
	u.fileSizeConfig = fileSizeConfig
}

// 初始化上传器
func (u *Uploader) RegisterUploader(uploader IUploader){
	u.uploader = uploader
}

// 设置md生成key
func (u *Uploader) SetKeyModel(mdKey bool){
	u.md = mdKey
}

func (u *Uploader) Upload(file *multipart.FileHeader, fileType FileType)(string, error){
	if size := u.fileSize(file.Size); !size {
		return "", errors.New("File exceeds size")
	}
	fileExt := path.Ext(file.Filename)
	if ok := u.fileSuffix(fileExt); !ok {
		return "", errors.New("Illegal file suffix")
	}

	// 如果是需要解压的文件类型，后缀必须是zip
	if fileType == UNZIP_FILE && fileExt != ".zip"{
		return "", errors.New("the suffix must be zip")
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	var b = make([]byte, file.Size)
	if _, err := f.Read(b); err != nil {
		return "", err
	}

	key := u.newKey(fileType, b, file.Filename, fileExt)
	data := bytes.NewBuffer(b)
	var fileData = FileData{
		Body: data,
		Size: file.Size,
	}

	return u.uploader.Upload(fileData, key, fileType)
}

func (u *Uploader)  newKey(fileType FileType, data []byte, fileName, fileExt string)(string){
	var typeDir string
	switch fileType {
	case IMAGE:
		//图片
		typeDir = "image"
	case VIDEO:
		//视频
		typeDir = "video"
	case AUDIO:
		//音频
		typeDir = "audio"
	default:
		//文件
		typeDir = "file"
	}

	var key = fmt.Sprintf("%s/%s", typeDir, fileName)
	if u.md {
		key = fmt.Sprintf("%s/", typeDir)
		key += hex.EncodeToString(u.Md5(data)) + fileExt
	}

	if fileType == UNZIP_FILE {
		key = fmt.Sprintf("%s", fileName)
	}

	return key
}

// MD5 md5 digest
func (u *Uploader)  Md5(plain []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(plain)
	cipher := md5Ctx.Sum(nil)
	return cipher
}

// MD5String md5 digest in string
func (u *Uploader)  md5String(plain string) string {
	cipher := u.Md5([]byte(plain))
	return hex.EncodeToString(cipher)
}

// FileSize 文件大小是否允许上传
func (u *Uploader)  fileSize(size int64) bool {
	//文件大小转换为MB
	fileSize := fmt.Sprintf("%.2f", float64(size)/float64(MB))
	if u.stringToFloat64(fileSize) > u.stringToFloat64(u.fileSizeConfig) {
		return false
	}
	return true
}

// FileSuffix 文件类型是否允许上传
func (u *Uploader)  fileSuffix(suffix string) bool {
	suffix = string([]byte(suffix)[1:len(suffix)])
	for _, v := range u.fileExtConfig {
		if v == suffix {
			return true
		}
	}
	return false
}

func (u *Uploader)  stringToFloat64(value string) float64 {
	data, err := strconv.ParseFloat(value, 64)
	if err != nil {
		zap.L().Info("StringToFloat64 error")
		return 0
	}
	return data
}