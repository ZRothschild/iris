// Package main展示了如何使用Iris内置的Jet视图引擎轻松使用嵌入在应用程序中的jet模板
// 此示例是https://github.com/CloudyKit/jet/tree/master/examples/asset_packaging的自定义分支，因此您可以并排注意到差异
// 例如，您不必在应用程序内部使用任何外部软件包，当Asset和AssetNames通过go-bindata之类的工具可用时，
// Iris可以手动为二进制数据构建模板加载器
//
// Package main shows how to use jet templates embedded in your application with ease using the Iris built-in Jet view engine.
// This example is a customized fork of https://github.com/CloudyKit/jet/tree/master/examples/asset_packaging, so you can
// notice the differences side by side. For example, you don't have to use any external package inside your application,
// Iris manually builds the template loader for binary data when Asset and AssetNames are available via tools like the go-bindata.
package main

import (
	"os"
	"strings"

	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	tmpl := iris.Jet("./views", ".jet").Binary(Asset, AssetNames)
	app.RegisterView(tmpl)

	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.jet")
	})

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	} else if !strings.HasPrefix(":", port) {
		port = ":" + port
	}

	app.Run(iris.Addr(port))
}
