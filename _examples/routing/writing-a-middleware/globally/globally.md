# IRIS全局路由函数注册
## 目录结构
> 主目录`globally`
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
	//将"before"处理程序注册为将要执行的第一个处理程序
	//在所有域的路由上。
	//或使用`UseGlobal`注册可跨子域触发的中间件。
	// app.Use(before)
	//将"after"处理程序注册为将要执行的最后一个处理程序
	//在所有域路由的处理程序之后。
	//
	//或使用`DoneGlobal`追加将在全局范围内触发的处理程序。
	// app.Done(after)

	// register the "before" handler as the first handler which will be executed
	// on all domain's routes.
	// Or use the `UseGlobal` to register a middleware which will fire across subdomains.
	// app.Use(before)
	// register the "after" handler as the last handler which will be executed
	// after all domain's routes' handler(s).
	//
	// Or use the `DoneGlobal` to append handlers that will be fired globally.
	// app.Done(after)

	//注册我们的路由。

	// register our routes.
	app.Get("/", indexHandler)
	app.Get("/contact", contactHandler)

	//这些调用的顺序无关紧要，`UseGlobal`和`DoneGlobal`
	//应用于现有路由与以后将要添加的路由
	//
	//请记住：`Use`和`Done`适用于当前路由组 与路由组的其子对象
	//因此，如果我们在路线注册之前使用`app.Use/Don`
	//在这种情况下，它将像UseGlobal/DoneGlobal 一样工作，因为`app`是初始路由组。

	// Order of those calls doesn't matter, `UseGlobal` and `DoneGlobal`
	// are applied to existing routes and future routes.
	//
	// Remember: the `Use` and `Done` are applied to the current party's and its children,
	// so if we used the `app.Use/Don`e before the routes registration
	// it would work like UseGlobal/DoneGlobal in this case, because the `app` is the root party.

	//有关更多信息，请参见`app.Party/PartyFunc`。

	// See `app.Party/PartyFunc` for more.
	app.UseGlobal(before)
	app.DoneGlobal(after)

	app.Run(iris.Addr(":8080"))
}

func before(ctx iris.Context) {
	shareInformation := "this is a sharable information between handlers"

	requestPath := ctx.Path()
	println("Before the indexHandler or contactHandler: " + requestPath)

	ctx.Values().Set("info", shareInformation)
	ctx.Next()
}

func after(ctx iris.Context) {
	println("After the indexHandler or contactHandler")
}

func indexHandler(ctx iris.Context) {
	println("Inside indexHandler")

	//从"before"处理程序中获取信息。

	// take the info from the "before" handler.
	info := ctx.Values().GetString("info")

	//向客户端写一些内容作为响应

	// write something to the client as a response.
	ctx.HTML("<h1>Response</h1>")
	ctx.HTML("<br/> Info: " + info)

	//执行通过`DoneGlobal`注册的"after"处理程序。
	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}

func contactHandler(ctx iris.Context) {
	println("Inside contactHandler")

	//向客户端写一些内容作为响应

	// write something to the client as a response.
	ctx.HTML("<h1>Contact</h1>")

	//执行通过`DoneGlobal`注册的"after"处理程序
	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}
```

## 介绍

1. 相当于构造函数与析构函数