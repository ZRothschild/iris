# ISIR DropzoneJS 文章 示例二
* [如何使用DropzoneJS和Go构建文件上传表单](https://hackernoon.com/how-to-build-a-file-upload-form-using-dropzonejs-and-go-8fb9f258a991)
* [如何使用DropzoneJS和Go在服务器上显示现有文件](https://hackernoon.com/how-to-display-existing-files-on-server-using-dropzonejs-and-go-53e24b57ba19)

# 内容
这是DropzoneJS + Go系列文章2之1

- [第1部分：如何构建文件上传表单](README-zh.md)
- [第2部分：如何在服务器上显示现有文件](README2-zh.md)
# DropzoneJS + Go：如何在服务器上显示现有文件
在本教程中，我们将向您展示在使用DropzoneJS和Go时如何在服务器上显示现有文件。本教程基于[如何使用DropzoneJS and Go构建文件上传表单](README-zh.md)。在继续阅读本教程中的内容之前，请确保已阅读它
## Table Of Content

- [预备](#预备)
- [修改服务器端](#修改服务器端)
- [修改客户端](#修改客户端)
- [参考资料](#参考资料)
- [结束](#结束)

## 预备

使用`go get github.com/nfnt/resize`安装go包`github.com/nfnt/resize`，我们需要它来创建缩略图

在前面的[教程](README-zh.md)中。我们已经建立了正确的DropzoneJs上传form表单。本教程不需要其他文件。 我们需要做的是对以下文件进行一些修改：

1. main.go
2. views/upload.html

让我们开始吧！

## 修改服务器端

在上一教程中。"/upload"所做的只是将上传的文件存储到服务器目录"./public/uploads"。因此，我们需要添加一段代码来检索存储文件的信息（名称和大小），并以JSON格式返回

将以下内容复制到"main.go"。阅读评论以获取详细信息

```golang
// main.go

package main

import (
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kataras/iris/v12"

	"github.com/nfnt/resize"
)

// $ go get -u github.com/nfnt/resize

const uploadsDir = "./public/uploads/"

type uploadedFile struct {
	// {name: "", size: } 是dropzone的唯一要求

	// {name: "", size: } are the dropzone's only requirements.
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type uploadedFiles struct {
	dir   string
	items []uploadedFile
	//切片是安全的，但是RWMutex对您来说是一个好习惯
	mu    sync.RWMutex // slices are safe but RWMutex is a good practise for you.
}

func scanUploads(dir string) *uploadedFiles {
	f := new(uploadedFiles)

	lindex := dir[len(dir)-1]
	if lindex != os.PathSeparator && lindex != '/' {
		dir += string(os.PathSeparator)
	}
	//根据需要创建目录，如果返回空的上传文件； 跳过扫描

	// create directories if necessary
	// and if, then return empty uploaded files; skipping the scan.
	if err := os.MkdirAll(dir, os.FileMode(0666)); err != nil {
		return f
	}
	//否则扫描给定的"dir" 以查找文件

	// otherwise scan the given "dir" for files.
	f.scan(dir)
	return f
}

func (f *uploadedFiles) scan(dir string) {
	f.dir = dir
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//如果是我们之前保存的目录或缩略图，请跳过它

		// if it's directory or a thumbnail we saved earlier, skip it.
		if info.IsDir() || strings.HasPrefix(info.Name(), "thumbnail_") {
			return nil
		}

		f.add(info.Name(), info.Size())
		return nil
	})
}

func (f *uploadedFiles) add(name string, size int64) uploadedFile {
	uf := uploadedFile{
		Name: name,
		Size: size,
	}

	f.mu.Lock()
	f.items = append(f.items, uf)
	f.mu.Unlock()

	return uf
}

func (f *uploadedFiles) createThumbnail(uf uploadedFile) {
	file, err := os.Open(path.Join(f.dir, uf.Name))
	if err != nil {
		return
	}
	defer file.Close()

	name := strings.ToLower(uf.Name)

	out, err := os.OpenFile(f.dir+"thumbnail_"+uf.Name,
		os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer out.Close()

	if strings.HasSuffix(name, ".jpg") {
		//将jpeg解码为image.Image

		// decode jpeg into image.Image
		img, err := jpeg.Decode(file)
		if err != nil {
			return
		}
		//将新图像写入文件

		// write new image to file
		resized := resize.Thumbnail(180, 180, img, resize.Lanczos3)
		jpeg.Encode(out, resized,
			&jpeg.Options{Quality: jpeg.DefaultQuality})

	} else if strings.HasSuffix(name, ".png") {
		img, err := png.Decode(file)
		if err != nil {
			return
		}
		//将新图像写入文件

		// write new image to file
		resized := resize.Thumbnail(180, 180, img, resize.Lanczos3) //速度较慢但分辨率更高 | slower but better res
		png.Encode(out, resized)
	}
	//依此类推...您明白了这一点，实际上，可以简化此代码
	
	// and so on... you got the point, this code can be simplify, as a practise.
}

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))

	app.HandleDir("/public", "./public")

	app.Get("/", func(ctx iris.Context) {
		ctx.View("upload.html")
	})

	files := scanUploads(uploadsDir)

	app.Get("/uploads", func(ctx iris.Context) {
		ctx.JSON(files.items)
	})

	app.Post("/upload", iris.LimitRequestBodySize(10<<20), func(ctx iris.Context) {
		//从dropzone请求中获取文件

		// Get the file from the dropzone request
		file, info, err := ctx.FormFile("file")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Application().Logger().Warnf("Error while uploading: %v", err.Error())
			return
		}

		defer file.Close()
		fname := info.Filename

		// Create a file with the same name
		// assuming that you have a folder named 'uploads'
		out, err := os.OpenFile(uploadsDir+fname,
			os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Application().Logger().Warnf("Error while preparing the new file: %v", err.Error())
			return
		}
		defer out.Close()

		io.Copy(out, file)
		//可选 将该文件添加到列表中，以便在刷新时可见

		// optionally, add that file to the list in order to be visible when refresh.
		uploadedFile := files.add(fname, info.Size)
		go files.createThumbnail(uploadedFile)
	})

	// start the server at http://localhost:8080
	app.Run(iris.Addr(":8080"))
}
```

## 修改客户端

将下面的内容复制到"./views/upload.html"。 我们将逐一进行修改
```html
<!-- /views/upload.html -->
<html>

<head>
    <title>DropzoneJS Uploader</title>

    <!-- 1 -->
    <link href="/public/css/dropzone.css" type="text/css" rel="stylesheet" />

    <!-- 2 -->
    <script src="/public/js/dropzone.js"></script>
    <!-- 4 -->
    <script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
    <!-- 5 -->
    <script>
        Dropzone.options.myDropzone = {
            paramName: "file", // The name that will be used to transfer the file
            init: function () {
                thisDropzone = this;
                // 6
                $.get('/uploads', function (data) {

                    if (data == null) {
                        return;
                    }
                    // 7
                    $.each(data, function (key, value) {
                        var mockFile = { name: value.name, size: value.size };

                        thisDropzone.emit("addedfile", mockFile);
                        thisDropzone.options.thumbnail.call(thisDropzone, mockFile, '/public/uploads/thumbnail_' + value.name);

                        // Make sure that there is no progress bar, etc...
                        thisDropzone.emit("complete", mockFile);
                    });

                });
            }
        };
    </script>
</head>

<body>

    <!-- 3 -->
    <form action="/upload" method="POST" class="dropzone" id="my-dropzone">
        <div class="fallback">
            <input name="file" type="file" multiple />
            <input type="submit" value="Upload" />
        </div>
    </form>
</body>

</html>
```
1.我们在页面中添加了Jquery库。实际上，这不是直接针对DropzoneJs的。我们仅使用Jquery的ajax函数**$.get**。您将在下面看到
2.我们在表单中添加了一个ID元素（my-dropzone）。这是必需的，因为我们需要将配置值传递给Dropzone。为此，我们必须有一个ID引用。这样我们就可以通过将值分配给Dropzone.options.myDropzone来配置它。配置Dropzone时，很多人会感到困惑。简单地说。不要将Dropzone用作Jquery插件，它具有自己的语法，您需要遵循它。
3.这开始了修改的主要部分。我们在这里所做的是传递一个函数来监听Dropzone的init事件。初始化Dropzone时将调用此事件。
4.通过ajax从新的"/uploads"检索文件详细信息。
5.使用服务器中的值创建模拟文件。mockFile只是具有名称和大小属性的JavaScript对象。然后，我们显式调用Dropzone的**addedfile**和**thumbnail**函数，以将现有文件放入Dropzone上传区域并生成其缩略图

### 运行服务器

在当前项目的文件夹中打开终端并执行：

```bash
$ go run main.go
Now listening on: http://localhost:8080
Application started. Press CTRL+C to shut down.
```
如果成功完成。 现在去上传一些图像并重新加载上传页面。 已上传的文件应自动显示在Dropzone区域中。

![上传的文件截图](with_files.png)

## 参考资料
- http://www.dropzonejs.com/#server-side-implementation
- https://www.startutorial.com/articles/view/how-to-build-a-file-upload-form-using-dropzonejs-and-php
- https://docs.iris-go.com
- https://github.com/kataras/iris/tree/master/_examples/tutorial/dropzonejs
## 结束
希望这个简单的教程对您的开发有所帮助。
如果您喜欢我的帖子，请在[Twitter](https://twitter.com/makismaropoulos)上关注我，并帮助宣传。 我需要您的支持才能继续。