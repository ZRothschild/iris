# IRIS路由前后置处理程序设置
## 目录结构
> 主目录`per-route`
```html
    —— main.go
```
## 代码示例
> `main.go`
```golang
package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()

	//或app.Use(before)和app.Done(after)

	// or app.Use(before) and app.Done(after).
	app.Get("/", before, mainHandler, after)

	// Use将一个中间件（前置处理程序）注册到所有当事方及其子级中
	//之后。
	//
	//`app`是根级，因此那些使用和完成的处理程序将在各处注册）

	// Use registers a middleware(prepend handlers) to all party's, and its children that will be registered
	// after.
	//
	// (`app` is the root children so those use and done handlers will be registered everywhere)
	app.Use(func(ctx iris.Context) {
		println(`before the party's routes and its children,
but this is not applied to the '/' route
because it's registered before the middleware, order matters.`)
		ctx.Next()
	})

	app.Done(func(ctx iris.Context) {
		println("this is executed always last, if the previous handler calls the `ctx.Next()`, it's the reverse of `.Use`")
		message := ctx.Values().GetString("message")
		println("message: " + message)
	})

	app.Get("/home", func(ctx iris.Context) {
		ctx.HTML("<h1> Home </h1>")
		ctx.Values().Set("message", "this is the home message, ip: "+ctx.RemoteAddr())
		ctx.Next() // call the done handlers.
	})

	child := app.Party("/child")
	child.Get("/", func(ctx iris.Context) {
		ctx.Writef(`this is the localhost:8080/child route.
All Use and Done handlers that are registered to the parent party,
are applied here as well.`)
		//调用完成的处理程序。
		ctx.Next() // call the done handlers.
	})

	app.Run(iris.Addr(":8080"))
}

func before(ctx iris.Context) {
	shareInformation := "this is a sharable information between handlers"

	requestPath := ctx.Path()
	println("Before the mainHandler: " + requestPath)

	ctx.Values().Set("info", shareInformation)
	//执行下一个处理程序，在本例中为主要处理程序。
	ctx.Next() // execute the next handler, in this case the main one.
}

func after(ctx iris.Context) {
	println("After the mainHandler")
}

func mainHandler(ctx iris.Context) {
	println("Inside mainHandler")

	//从"before"处理程序中获取信息

	// take the info from the "before" handler.
	info := ctx.Values().GetString("info")

	//向客户端写一些内容作为响应。

	// write something to the client as a response.
	ctx.HTML("<h1>Response</h1>")
	ctx.HTML("<br/> Info: " + info)

	//执行"after".
	ctx.Next() // execute the "after".
}
```