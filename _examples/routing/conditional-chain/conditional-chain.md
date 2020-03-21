# IRIS条件链路由
## 目录结构
> 主目录`conditional-chain`
```html
    —— main.go
    —— main_test.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"github.com/kataras/iris/v12"
)

func newApp() *iris.Application {
	app := iris.New()
	v1 := app.Party("/api/v1")

	myFilter := func(ctx iris.Context) bool {
		//请勿在生产环境中执行此操作，请使用会话或/和数据库调用等。

		// don't do that on production, use session or/and database calls and etc.
		ok, _ := ctx.URLParamBool("admin")
		return ok
	}

	onlyWhenFilter1 := func(ctx iris.Context) {
		ctx.Application().Logger().Infof("admin: %#+v", ctx.URLParams())
		ctx.Writef("<title>Admin</title>\n")
		ctx.Next()
	}

	onlyWhenFilter2 := func(ctx iris.Context) {
		//您可以一直使用上一个求情处理方法存储的数据
		//执行类似ofc的操作。
		//
		//当前路由处理方法设置 ：ctx.Values().Set("is_admin", true)
		//下一个路由处理方法可以获取到上一个设置的值 ：isAdmin := ctx.Values().GetBoolDefault("is_admin", false)
		//
		//，但让我们简化一下：

		// You can always use the per-request storage
		// to perform actions like this ofc.
		//
		// this handler: ctx.Values().Set("is_admin", true)
		// next handler: isAdmin := ctx.Values().GetBoolDefault("is_admin", false)
		//
		// but, let's simplify it:
		ctx.HTML("<h1>Hello Admin</h1><br>")
		ctx.Next()
	}

	// 这里：
	//它可以在任何地方注册为中间件。
	//它将触发`onlyWhenFilter1`和`onlyWhenFilter2`作为中间件（使用ctx.Next（））
	//如果myFilter通过，否则它将通过忽略ctx.Next（）继续处理程序链
	//`onlyWhenFilter1`和`onlyWhenFilter2`。

	// HERE:
	// It can be registered anywhere, as a middleware.
	// It will fire the `onlyWhenFilter1` and `onlyWhenFilter2` as middlewares (with ctx.Next())
	// if myFilter pass otherwise it will just continue the handler chain with ctx.Next() by ignoring
	// the `onlyWhenFilter1` and `onlyWhenFilter2`.
	myMiddleware := iris.NewConditionalHandler(myFilter, onlyWhenFilter1, onlyWhenFilter2)

	v1UsersRouter := v1.Party("/users", myMiddleware)
	v1UsersRouter.Get("/", func(ctx iris.Context) {
		ctx.HTML("requested: <b>/api/v1/users</b>")
	})

	return app
}

func main() {
	app := newApp()

	// http://localhost:8080/api/v1/users
	// http://localhost:8080/api/v1/users?admin=true
	app.Run(iris.Addr(":8080"))
}
```
> `main_test.go`

```go
package main

import (
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func TestNewConditionalHandler(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	e.GET("/api/v1/users").Expect().Status(httptest.StatusOK).
		Body().Equal("requested: <b>/api/v1/users</b>")
	e.GET("/api/v1/users").WithQuery("admin", "true").Expect().Status(httptest.StatusOK).
		Body().Equal("<title>Admin</title>\n<h1>Hello Admin</h1><br>requested: <b>/api/v1/users</b>")
}
```