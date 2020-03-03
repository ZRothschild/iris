# `route` 自定义高级的路由
## 目录结构
> 主目录`custom-high-level-router`
```html
    —— main.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

/*
	路由器应包含以下三种方法：
   -HandleRequest应该基于上下文处理请求。
		HandleRequest(ctx iris.Context)
    -构建路由处理函数，它是在路由器的BuildRouter上调用的。
		Build(provider router.RoutesProvider) error
    -RouteExists报告是否存在特定路由。
       RouteExists(ctx iris.Context, method, path string) bool

有关更详细，完整和有用的示例
您可以查看：
https://github.com/kataras/iris/tree/master/core/router/handler.go
这就完成了这个接口`router＃RequestHandler`。
*/

/* A Router should contain all three of the following methods:
   - HandleRequest should handle the request based on the Context.
	  HandleRequest(ctx iris.Context)
   - Build should builds the handler, it's being called on router's BuildRouter.
	  Build(provider router.RoutesProvider) error
   - RouteExists reports whether a particular route exists.
      RouteExists(ctx iris.Context, method, path string) bool

For a more detailed, complete and useful example
you can take a look at the iris' router itself which is located at:
https://github.com/kataras/iris/tree/master/core/router/handler.go
which completes this exact interface, the `router#RequestHandler`.
*/
type customRouter struct {
	//复制一个路由（这比较安全，因为如果没有`app.RefreshRouter`调用，您将无法在服务时间更改路由）：
	// [] router.Route
	//或只是期望整个路由提供者：

	// a copy of routes (safer because you will not be able to alter a route on serve-time without a `app.RefreshRouter` call):
	// []router.Route
	// or just expect the whole routes provider:
	provider router.RoutesProvider
}
// HandleRequest一个愚蠢的示例，该示例仅基于请求路径的第一部分来查找路由
//它也必须是静态的，其余的将填充参数。

// HandleRequest a silly example which finds routes based only on the first part of the requested path
// which must be a static one as well, the rest goes to fill the parameters.
func (r *customRouter) HandleRequest(ctx iris.Context) {
	path := ctx.Path()
	ctx.Application().Logger().Infof("Requested resource path: %s", path)

	parts := strings.Split(path, "/")[1:]
	staticPath := "/" + parts[0]
	for _, route := range r.provider.GetRoutes() {
		if strings.HasPrefix(route.Path, staticPath) && route.Method == ctx.Method() {
			paramParts := parts[1:]
			for _, paramValue := range paramParts {
				for _, p := range route.Tmpl().Params {
					ctx.Params().Set(p.Name, paramValue)
				}
			}

			ctx.SetCurrentRouteName(route.Name)
			ctx.Do(route.Handlers)
			return
		}
	}
	//如果没有找到...
	// if nothing found...
	ctx.StatusCode(iris.StatusNotFound)
}

func (r *customRouter) Build(provider router.RoutesProvider) error {
	for _, route := range provider.GetRoutes() {
		//根据您的自定义逻辑进行任何必要的验证或对话
		//，但始终为每个已注册的路由运行route.BuildHandlers()

		// do any necessary validation or conversations based on your custom logic here
		// but always run the "BuildHandlers" for each registered route.
		route.BuildHandlers()
		// [...] r.routes = append(r.routes, *route)
	}

	r.provider = provider
	return nil
}

func (r *customRouter) RouteExists(ctx iris.Context, method, path string) bool {
	// [...]
	return false
}

func main() {
	app := iris.New()

	//如果您想知道，参数类型和宏 例如"{param:string $func()}"仍然可以在内部使用
	//您的自定义路由器（通过路由的处理程序获取）
	//因为它们是设置的中间件钩子，所以您不必实现手动处理它们的逻辑，
	//尽管您必须匹配请求的路径是什么路由并填充ctx.Params()，但这是自定义路由器的工作。

	// In case you are wondering, the parameter types and macros like "{param:string $func()}" still work inside
	// your custom router if you fetch by the Route's Handler
	// because they are middlewares under the hood, so you don't have to implement the logic of handling them manually,
	// though you have to match what requested path is what route and fill the ctx.Params(), this is the work of your custom router.
	app.Get("/hello/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef("Hello %s\n", name)
	})

	app.Get("/cs/{num:uint64 min(10) else 400}", func(ctx iris.Context) {
		num := ctx.Params().GetUint64Default("num", 0)
		ctx.Writef("num is: %d\n", num)
	})

	//通过使用iris/context.Context将现有router替换为定制的router
	//您必须在app.Run之前和注册路由之后使用app.BuildRouter方法。
	//您应该将自定义router的实例作为第二个输入参数传递，该输入必须完成`router＃RequestHandler`
	//接口，如上所示。
	//
	//要了解如何在没有直接Iris上下文支持的情况下构建更底层的内容（您也可以手动执行此操作）
	//可以参考"custom-wrapper"示例。
	
	// To replace the existing router with a customized one by using the iris/context.Context
	// you have to use the `app.BuildRouter` method before `app.Run` and after the routes registered.
	// You should pass your custom router's instance as the second input arg, which must completes the `router#RequestHandler`
	// interface as shown above.
	//
	// To see how you can build something even more low-level without direct iris' context support (you can do that manually as well)
	// navigate to the "custom-wrapper" example instead.
	myCustomRouter := new(customRouter)
	app.BuildRouter(app.ContextPool, myCustomRouter, app.APIBuilder, true)

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
```