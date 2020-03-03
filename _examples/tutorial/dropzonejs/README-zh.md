# 文章
* [如何使用DropzoneJS和Go构建文件上传表单](https://hackernoon.com/how-to-build-a-file-upload-form-using-dropzonejs-and-go-8fb9f258a991)
* [如何使用DropzoneJS和Go在服务器上显示现有文件](https://hackernoon.com/how-to-display-existing-files-on-server-using-dropzonejs-and-go-53e24b57ba19)

# 内容
这是DropzoneJS + Go系列文章2之1

- [第1部分：如何构建文件上传表单](README-zh.md)
- [第2部分：如何在服务器上显示现有文件](README2-zh.md)

# DropzoneJS + Go：如何构建文件上传表单
[DropzoneJS](https://github.com/enyo/dropzone) 是一个开放源代码库，提供带有图像预览的拖放文件上传。 这是一个很棒的JavaScript库，
实际上甚至不依赖JQuery。在本教程中，我们将使用DropzoneJS构建一个多文件上传表单，后端将由Go [Iris](https://iris-go.com).
## 表中的内容
- [准备](#准备)
- [dropzone.js使用](#dropzone.js使用)
- [Go使用](#go使用)
## 准备
1. 下载 [Go(Golang)](https://golang.org/dl), 如图所示设置计算机，然后继续执行2
2. 安装 [Iris](https://github.com/kataras/iris); 打开一个终端并执行`go get -u github.com/kataras/iris`
3. 从这里下载DropzoneJS从 [this URL](https://raw.githubusercontent.com/enyo/dropzone/master/dist/dropzone.js). DropzoneJS不依赖JQuery，您不必担心，升级JQuery版本会破坏您的应用程序
4. 从这里下载dropzone.css从 [this URL](https://raw.githubusercontent.com/enyo/dropzone/master/dist/dropzone.css), 如果您想要一些已经制作好的CSS
5. 创建一个文件夹"./public/uploads"，该文件夹用于存储上传的文件
6. 创建一个文件"./views/upload.html"，该文件用于首页
7. 创建一个文件"./main.go"，该文件用于处理后端文件上传过程
准备之后，您的文件夹和文件结构应如下所示：
![文件夹和文件结构](folder_structure.png)
## dropzone.js使用
打开文件"./views/upload.html"，让我们创建一个DropzoneJs表单

将下面的内容复制到"./views/upload.html"，我们将逐一逐行检查代码
```html
<!-- /views/upload.html -->
<html>

<head>
    <title>DropzoneJS Uploader</title>

    <!-- 1 -->
    <link href="/public/css/dropzone.css" type="text/css" rel="stylesheet" />

    <!-- 2 -->
    <script src="/public/js/dropzone.js"></script>
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
1. 包括CSS样式表
2. 包括DropzoneJS JavaScript库
3. 使用CSS类"dropzone"创建一个上传表单，并且"action"是路由路径"/upload"。请注意，我们确实为后备模式创建了一个输入字段。这全部由DropzoneJS库本身处理。我们需要做的就是将CSS类"dropzone"分配给表单。 默认情况下，DropzoneJS将查找所有带有"dropzone"类的表单，并自动将其自身附加到该表单。
## go使用
现在，您已经进入了教程的最后一部分。 在本节中，我们会将从DropzoneJS发送的文件存储到"./public/uploads"文件夹中。

打开"main.go" 并复制以下代码：
```go
// main.go

package main

import (
    "os"
    "io"
    "strings"

    "github.com/kataras/iris/v12"
)

const uploadsDir = "./public/uploads/"

func main() {
    app := iris.New()

    //注册模板
    app.RegisterView(iris.HTML("./views", ".html"))
    //设置/public路由路径以静态服务./public / ...内容

    // Make the /public route path to statically serve the ./public/... contents
    app.HandleDir("/public", "./public")
    //渲染实际形式

    // Render the actual form
    // GET: http://localhost:8080
    app.Get("/", func(ctx iris.Context) {
        ctx.View("upload.html")
    })
    //将文件上传到服务器

    // Upload the file to the server
    // POST: http://localhost:8080/upload
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
        //创建一个同名文件
        //假设您有一个名为'uploads'的文件夹

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
    })
    //启动服务器 http://localhost:8080

    // Start the server at http://localhost:8080
    app.Run(iris.Addr(":8080"))
}
```

1.创建一个新的Iris应用。
2.从"views"文件夹中注册并加载模板。
3.设置"/public"路由路径以静态服务../public/...文件夹的内容
4.创建一路由来服务上传的form表单页面
5.创建一条路由来处理DropzoneJS表单中的POST表单数据
6.为目标文件夹声明一个变量。
7.如果将文件发送到页面，则将文件对象存储到一个临时的“file”变量中
8.根据uploadsDir +上传的文件名将上传的文件移动到目标位置
### 运行服务器
在当前项目的文件夹中打开终端并执行：
```bash
$ go run main.go
Now listening on: http://localhost:8080
Application started. Press CTRL+C to shut down.
```

现在转到浏览器，并导航到http://localhost:8080，您应该能够看到如下页面：

![没有文件截图](no_files.png)
![上传的文件截图](with_files.png)