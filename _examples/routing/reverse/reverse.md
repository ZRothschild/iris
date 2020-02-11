# `route`路由转换
## 目录结构
> 主目录`reverse`
```html
    —— main.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

func main() {
	app := iris.New()
	//需要在视图引擎外部进行手动反向路由。
	// {{ urlpath "routename" "path" "values" "here"}}通常不需要

	// need for manually reverse routing when needed outside of view engine.
	// you normally don't need it because of the {{ urlpath "routename" "path" "values" "here"}}
	rv := router.NewRoutePathReverser(app)

	myroute := app.Get("/anything/{anythingparameter:path}", func(ctx iris.Context) {
		paramValue := ctx.Params().Get("anythingparameter")
		ctx.Writef("The path after /anything is: %s", paramValue)
	})

	myroute.Name = "myroute"
	//对链接很有用，尽管iris的视图引擎已经具有{{ urlpath "routename" "path values"}}

	// useful for links, although iris' view engine has the {{ urlpath "routename" "path values"}} already.
	app.Get("/reverse_myroute", func(ctx iris.Context) {
		myrouteRequestPath := rv.Path(myroute.Name, "any/path")
		ctx.HTML("Should be <b>/anything/any/path</b>: " + myrouteRequestPath)
	})
	//执行一条路由，类似于重定向但没有重定向:)

	// execute a route, similar to redirect but without redirect :)
	app.Get("/execute_myroute", func(ctx iris.Context) {
		//就像客户端调用它一样
		ctx.Exec("GET", "/anything/any/path") // like it was called by the client.
	})

	// http://localhost:8080/reverse_myroute
	// http://localhost:8080/execute_myroute
	// http://localhost:8080/anything/any/path/here

	//有关更多反向路由示例，请参见view/template_html_4示例
	//使用反向路由器组件以及{{url}}和{{urlpath}}模板函数。

	// See view/template_html_4 example for more reverse routing examples
	// using the reverse router component and the {{url}} and {{urlpath}} template functions.
	app.Run(iris.Addr(":8080"))
}
```

### 介绍

1. 相当于执行当前路由，再去指定执行一个路由类似go to