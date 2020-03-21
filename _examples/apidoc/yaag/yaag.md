# `YAAG：golang web API` 文档生成器
## YAAG介绍
`Golang`非常适合开发`Web`应用程序。人们已经创建了许多优秀的`Web`框架，`Web`帮助程序库。 
如果我们考虑`Golang`中的整个`Web`应用程序生态系统，似乎很完整了，但我们却少了一个编写`API`文档的中间件。
因此，我们为基于`Golang`的`Web`应用程序创建了第一个`API`文档生成器。

大多数`Web`服务都将其`API`暴露给移动或第三方开发人员。 记录它们有点痛苦。 
我们正在努力减轻痛苦，至少对于您不必向世界公开您的文档的内部项目。 
`YAAG`生成简单的基于引导程序的`API`文档，无需编写任何注释。

`YAAG`是一个中间件。 您必须在路线中添加`YAAG`处理程序，您就完成了。 
继续使用`POSTMAN`，`Curl`或任何客户端调用您的`API`，`YAAG`将继续更新`API Doc html`。
注意：我们还生成一个包含所有`API`调用数据的`json`文件）

## `YAAG `生成`iris web`框架项目`API`文档
### API 生成步骤
- 下载`YAAG`中间件
> go get github.com/betacraft/yaag/...
- 导入依赖包
> import github.com/betacraft/yaag/yaag
> Import github.com/betacraft/yaag/irisyaag
- 初始化`yaag`
> yaag.Init(&yaag.Config(On: true, DocTile: "Iris", DocPath: "apidoc.html"))
- 注册`yaag`中间件
> app.Use(irisyaag.New())
> irisyaag记录响应主体并向apidoc提供所有必要的信息
### 目录结构
> 主目录`iris`

```html
    —— apidoc.html (执行命令后生成)
    —— apidoc.html.json (执行命令后生成)
    —— main.go
```
## 代码示例
> `main.go`
```golang
package main

import (
	"github.com/kataras/iris/v12"

	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
)

type myXML struct {
	Result string `xml:"result"`
}

func main() {
	app := iris.New()

	yaag.Init(&yaag.Config{ // 重要说明，初始化中间件 | IMPORTANT, init the middleware.
		On:       true,		//是否开启自动生成API文档功能
		DocTitle: "Iris",
		DocPath:  "apidoc.html",
		BaseUrls: map[string]string{"Production": "", "Staging": ""},
	})
	app.Use(irisyaag.New()) // <- 重要，注册中间件 | IMPORTANT, register the middleware.

	app.Get("/json", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"result": "Hello World!"})
	})

	app.Get("/plain", func(ctx iris.Context) {
		ctx.Text("Hello World!")
	})

	app.Get("/xml", func(ctx iris.Context) {
		ctx.XML(myXML{Result: "Hello World!"})
	})

	app.Get("/complex", func(ctx iris.Context) {
		value := ctx.URLParam("key")
		ctx.JSON(iris.Map{"value": value})
	})

	//运行我们的HTTP服务器。
	//
	//“yaag”的文档没有说明以下内容，但是在Iris中，我们在为您提供的内容方面非常谨慎。
	//每个传入的请求都会重新生成和更新“apidoc.html”文件。
	//建议：
	//编写调用这些处理程序的测试，保存生成的“apidoc.html”。
	//在生产中关闭yaag中间件。
	//
	//用法示例：
	//访问所有路径并打开生成的“apidoc.html”文件，以查看API的自动文档。

	// Run our HTTP Server.
	//
	// Documentation of "yaag" doesn't note the follow, but in Iris we are careful on what
	// we provide to you.
	//
	// Each incoming request results on re-generation and update of the "apidoc.html" file.
	// Recommentation:
	// Write tests that calls those handlers, save the generated "apidoc.html".
	// Turn off the yaag middleware when in production.
	//
	// Example usage:
	// Visit all paths and open the generated "apidoc.html" file to see the API's automated docs.
	app.Run(iris.Addr(":8080"))
}
```
### 提示
1. 运行上面的例子，并请求其中任意一个接口，会生成`apidoc.html`，`apidoc.html.json`两个文件
2. 如果不关闭`yaag`中间件，`apidoc.html`，`apidoc.html.json`文件会随着每一次请求而重新生成
3. 如果没有翻墙，可能看不到生成效果，因为`apidoc.html`文件的引入了很多谷歌的`cnd`,替换成功国内支持即可