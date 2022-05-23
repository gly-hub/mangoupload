### mangoupload文件上传集成

### 集成

* local 文件夹
* ali **OOS**
* huawei **OBS**
* qiniu **QBox**
* tencent **COS**

### 使用方法

#### 1.初始化上传器

```golang
global.Uploader = mango_upload.NewUploader()
global.Uploader.SetKeyModel(true)

global.Uploader.SetUploadLimit(config.Conf.FileSet.FileSize, config.Conf.FileSet.FileExt)
```

#### 2.初始化指定上传器

（以本地上传器为例）

```golang
localUploader := mango_upload.NewLocalUploader(mango_upload.LocalConfig{
    Domain:    config.Conf.FileLocal.Domain,
    LocalPath: config.Conf.FileLocal.Path,
    RoutePath: "v1/local",
})

global.Uploader.RegisterUploader(localUploader)
```

#### 3.注册上传器

```golang
global.LocalUploader.RegisterUploader(localUploader)
```

#### 4.上传示例

```golang
func UploadFile(c *gin.Context) {
	fileType := c.PostForm("type")
	file, err := c.FormFile("upload_file")

	switch fileType {
	case "1":
		//图片
		fileType = "image"
	case "2":
		//视频
		fileType = "video"
	case "3":
		//音频
		fileType = "audio"
	case "4":
		fileType = "unzipFile"
	default:
		//文件
		fileType = "file"
	}

	tools.HasError(err, "", -1)

	url, err := global.Uploader.Upload(file, mango_upload.FileType(fileType))
	tools.HasError(err, "", -1)
	var r = map[string]string{
		"url": url,
	}
	app.OK(c, r, "")
}
```

