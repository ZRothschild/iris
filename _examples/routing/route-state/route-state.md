# `route`路由状态设置
## 目录结构
> 主目录`route-state`
```html
    —— main.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	none := app.None("/invisible/{username}", func(ctx iris.Context) {
		ctx.Writef("Hello %s with method: %s", ctx.Params().Get("username"), ctx.Method())

		if from := ctx.Values().GetString("from"); from != "" {
			ctx.Writef("\nI see that you're coming from %s", from)
		}
	})

	app.Get("/change", func(ctx iris.Context) {
		if none.IsOnline() {
			none.Method = iris.MethodNone
		} else {
			none.Method = iris.MethodGet
		}
		// refresh在服务时重新构建路由器，以便收到有关其新路由的通知。

		// refresh re-builds the router at serve-time in order to be notified for its new routes.
		app.RefreshRouter()
	})

	app.Get("/execute", func(ctx iris.Context) {
		if !none.IsOnline() {
			ctx.Values().Set("from", "/execute with offline access")
			ctx.Exec("NONE", "/invisible/iris")
			return
		}
		//与调用/change并更改路由状态时导航到"http://localhost:8080/invisible/iris"相同
		//从"offline" 到 "online"

		// same as navigating to "http://localhost:8080/invisible/iris" when /change has being invoked and route state changed
		// from "offline" to "online"

		//从外部上下文调用Exec时，可以共享值和会话
		ctx.Values().Set("from", "/execute") // values and session can be shared when calling Exec from a "foreign" context.
		// 	ctx.Exec("NONE", "/invisible/iris")
		// or after "/change":
		ctx.Exec("GET", "/invisible/iris")
	})

	app.Run(iris.Addr(":8080"))
}
```

## 介绍

1. 设置路由状态可以跳过不需要调用的处理函数