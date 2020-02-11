package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	//现在可以正常工作了，
	//将处理所有GET请求
	// 除了：
	// /                     -> 因为 app.Get("/", ...)
	// /other/anything/here  -> 因为 app.Get("/other/{paramother:path}", ...)
	// /other2/anything/here -> 因为 app.Get("/other2/{paramothersecond:path}", ...)
	// /other2/static2       -> 因为 app.Get("/other2/static", ...)
	//
	//它与其余路由没有冲突，没有路由性能成本！
	//
	//即 /something/here/that/cannot/be/found/by/other/registered/routes/order/not/matters

	// this works as expected now,
	// will handle all GET requests
	// except:
	// /                     -> because of app.Get("/", ...)
	// /other/anything/here  -> because of app.Get("/other/{paramother:path}", ...)
	// /other2/anything/here -> because of app.Get("/other2/{paramothersecond:path}", ...)
	// /other2/static2        -> because of app.Get("/other2/static", ...)
	//
	// It isn't conflicts with the rest of the routes, without routing performance cost!
	//
	// i.e /something/here/that/cannot/be/found/by/other/registered/routes/order/not/matters
	app.Get("/{p:path}", h)
	// app.Get("/static/{p:path}", staticWildcardH)

	//这只会处理 GET /
	// this will handle only GET /
	app.Get("/", staticPath)

	//这将处理所有以"/other/"开头的GET请求
	// this will handle all GET requests starting with "/other/"
	//
	// i.e /other/more/than/one/path/parts
	app.Get("/other/{paramother:path}", other)

	//这将处理所有以"/other2/"开头的GET请求
	// /other2/static 除外（由于下一条静态路由）

	// this will handle all GET requests starting with "/other2/"
	// except /other2/static (because of the next static route)
	//
	// i.e /other2/more/than/one/path/parts
	app.Get("/other2/{paramothersecond:path}", other2)
	//这只会处理GET "/other2/static"
	// this will handle only GET "/other2/static"
	app.Get("/other2/static2", staticPathOther2)

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func h(ctx iris.Context) {
	param := ctx.Params().Get("p")
	ctx.WriteString(param)
}

func staticWildcardH(ctx iris.Context) {
	param := ctx.Params().Get("p")
	ctx.WriteString("from staticWildcardH: param=" + param)
}

func other(ctx iris.Context) {
	param := ctx.Params().Get("paramother")
	ctx.Writef("from other: %s", param)
}

func other2(ctx iris.Context) {
	param := ctx.Params().Get("paramothersecond")
	ctx.Writef("from other2: %s", param)
}

func staticPath(ctx iris.Context) {
	ctx.Writef("from the static path(/): %s", ctx.Path())
}

func staticPathOther2(ctx iris.Context) {
	ctx.Writef("from the static path(/other2/static2): %s", ctx.Path())
}
