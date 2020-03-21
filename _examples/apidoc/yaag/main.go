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
