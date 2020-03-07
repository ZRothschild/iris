package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()

	tmpl := iris.Pug("./templates", ".pug")
	//根据每个请求重新加载模板（开发模式）
	tmpl.Reload(true) // reload templates on each request (development mode)
	//在此处添加模板功能
	tmpl.AddFunc("greet", func(s string) string { // add your template func here.
		return "Greetings " + s + "!"
	})

	app.RegisterView(tmpl)

	app.Get("/", index)

	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}

func index(ctx iris.Context) {
	ctx.ViewData("pageTitle", "My Index Page")
	ctx.ViewData("youAreUsingJade", true)
	//问：为什么需要扩展名.pug？
	//答：因为您可以为每个Iris应用程序注册多个视图引擎

	// Q: why need extension .pug?
	// A: Because you can register more than one view engine per Iris application.
	ctx.View("index.pug")
}
