package main

import (
	"github.com/kataras/iris/v12"
)

type mypage struct {
	Title   string
	Message string
}

func main() {
	app := iris.New()

	app.RegisterView(iris.HTML("./templates", ".html").Layout("layout.html"))
	//提示：追加..Reload(true)以在每个请求上重新加载模板

	// TIP: append .Reload(true) to reload the templates on each request.

	app.Get("/", func(ctx iris.Context) {
		ctx.Gzip(true)
		ctx.ViewData("", mypage{"My Page title", "Hello world!"})
		ctx.View("mypage.html")

		//注意：您可以传递"layout" : "otherLayout.html"绕过配置的Layout属性
		//或view.NoLayout禁用此渲染操作的布局。第三个是可选参数

		// Note that: you can pass "layout" : "otherLayout.html" to bypass the config's Layout property
		// or view.NoLayout to disable layout on this render action.
		// third is an optional parameter
	})

	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}
