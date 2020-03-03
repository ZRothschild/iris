package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	/*
	 * Setup static files 设置静态文件
	 */

	app.HandleDir("/assets", "./public/assets")
	app.HandleDir("/upload_resources", "./public/upload_resources")

	dashboard := app.Party("dashboard.")
	{
		dashboard.Get("/", func(ctx iris.Context) {
			ctx.Writef("HEY FROM dashboard")
		})
	}
	system := app.Party("system.")
	{
		system.Get("/", func(ctx iris.Context) {
			ctx.Writef("HEY FROM system")
		})
	}

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("HEY FROM frontend /")
	})
	// http://domain.local:80
	// http://dashboard.local
	// http://system.local

	//确保您在浏览器中添加"http"，因为.local是一个虚拟域，
	// 我们认为在这种情况下，您可以将任何语法正确的名称声明为iris的子域。

	// Make sure you prepend the "http" in your browser
	// because .local is a virtual domain we think to show case you
	// that you can declare any syntactical correct name as a subdomain in iris.

	//对于初学者：查看../hosts文件
	app.Run(iris.Addr("domain.local:80")) // for beginners: look ../hosts file
}
