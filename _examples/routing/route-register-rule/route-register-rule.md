# IRIS路由注册规则
## 目录结构
> 主目录`route-register-rule`
```html
    —— main.go
    —— main_test.go
```
## 代码示例
> `main.go`
```golang
package main

import "github.com/kataras/iris/v12"

func main() {
	app := newApp()
	//浏览https://github.com/kataras/iris/issues/1448了解详细信息。

	// Navigate through https://github.com/kataras/iris/issues/1448 for details.
	//
	// GET: http://localhost:8080
	// POST, PUT, DELETE, CONNECT, HEAD, PATCH, OPTIONS, TRACE : http://localhost:8080
	app.Listen(":8080")
}

func newApp() *iris.Application {
	app := iris.New()
	//跳过并且不要覆盖现有的已转接路由，请正常继续。
	//适用于某个方及其子级，在此情况下，适用于整个应用程序的路由。

	// Skip and do NOT override existing regitered route, continue normally.
	// Applies to a Party and its children, in this case the whole application's routes.
	app.SetRegisterRule(iris.RouteSkip)

	/* Read also:

	//默认行为是将对`app.Any`调用中的anyHandler覆盖getHandler。
	app.SetRegistRule（iris.RouteOverride）

	//停止执行并在服务器引导之前引发错误。
	app.SetRegisterRule（iris.RouteError）

	// The default behavior, will override the getHandler to anyHandler on `app.Any` call.
	app.SetRegistRule(iris.RouteOverride)

	// Stops the execution and fires an error before server boot.
	app.SetRegisterRule(iris.RouteError)
	*/

	app.Get("/", getHandler)
	// app.Any因为`iris.RouteSkip`规则而没有覆盖之前的GET路由。

	// app.Any does NOT override the previous GET route because of `iris.RouteSkip` rule.
	app.Any("/", anyHandler)

	return app
}

func getHandler(ctx iris.Context) {
	ctx.Writef("From get handle %s", ctx.GetCurrentRoute().Trace())
}

func anyHandler(ctx iris.Context) {
	ctx.Writef("From %s", ctx.GetCurrentRoute().Trace())
}
```

> `main_test.go`
```golang
package main

import (
	"testing"

	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/httptest"
)

func TestRouteRegisterRuleExample(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	for _, method := range router.AllMethods {
		tt := e.Request(method, "/").Expect().Status(httptest.StatusOK).Body()
		if method == "GET" {
			tt.Equal("From [./main.go:28] GET: / -> github.com/kataras/iris/v12/_examples/routing/route-register-rule.getHandler()")
		} else {
			tt.Equal("From [./main.go:30] " + method + ": / -> github.com/kataras/iris/v12/_examples/routing/route-register-rule.anyHandler()")
		}
	}
}
```