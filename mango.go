/**
* @Author vangogh
* @Description 文件类型
* @File:  filetype
* @Datetime 2022/4/20 10:38
**/
package mango_upload

type FileType string

const(
	IMAGE FileType = "image" // 图片文件
	VIDEO FileType = "video" // 视频
	AUDIO FileType = "audio" // 音频
	FILE FileType = "file" // 普通文件
	UNZIP_FILE FileType = "unzipFile" // 压缩文件，该类型表明文件需要解压
	APK_FILE FileType = "apkFile" // 安卓包文件
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

var (
	FileSizeConfig = "10" // 单位M
	FileExtConfig = []string{"jpeg", "jpg", "png", "mp4", "mp3", "gif", "bmp", "zip"}
)
