# `route`基本使用，基础用法
## 目录结构
> 主目录`basic`
```html
    —— main.go
    —— test_main.go
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
	app.Logger().SetLevel("debug")

	//注册一个自定处理处理 hhttp 状态为404未找到路径错误的路由处理函数
	//触发条件是路由没有被找到，或通过手动调用 ctx.StatusCode(iris.StatusNotFound)

	// registers a custom handler for 404 not found http (error) status code,
	// fires when route not found or manually by ctx.StatusCode(iris.StatusNotFound).
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)

	// GET -> HTTP 的请求方法
	// / -> /路由名称
	// func(ctx iris.Context) -> 路由处理函数
	//
	//第三个可变参数应该包含一个或多个路由处理函数，他们将被顺序执行
	//例子如下:
	
	// GET -> HTTP Method
	// / -> Path
	// func(ctx iris.Context) -> The route's handler.
	//
	// Third receiver should contains the route's handler(s), they are executed by order.
	app.Handle("GET", "/", func(ctx iris.Context) {

		// 可以参考 https://github.com/kataras/iris/wiki/Routing-context-methods
		//详细介绍了所有 context 的所有可用方法不仅仅是 ctx.Path()
		
		// navigate to the https://github.com/kataras/iris/wiki/Routing-context-methods
		// to overview all context's method.
		ctx.HTML("Hello from " + ctx.Path()) // Hello from /
	})

	app.Get("/home", func(ctx iris.Context) {
		ctx.Writef(`Same as app.Handle("GET", "/", [...])`)
	})
	
	//同一路径中的不同路径参数类型。
	
	// Different path parameters types in the same path.
	app.Get("/u/{p:path}", func(ctx iris.Context) {
		ctx.Writef(":string, :int, :uint, :alphabetical and :path in the same path pattern.")
	})

	app.Get("/u/{username:string}", func(ctx iris.Context) {
		ctx.Writef("before username (string), current route name: %s\n", ctx.RouteName())
		ctx.Next()
	}, func(ctx iris.Context) {
		ctx.Writef("username (string): %s", ctx.Params().Get("username"))
	})

	app.Get("/u/{id:int}", func(ctx iris.Context) {
		ctx.Writef("before id (int), current route name: %s\n", ctx.RouteName())
		ctx.Next()
	}, func(ctx iris.Context) {
		ctx.Writef("id (int): %d", ctx.Params().GetIntDefault("id", 0))
	})

	app.Get("/u/{uid:uint}", func(ctx iris.Context) {
		ctx.Writef("before uid (uint), current route name: %s\n", ctx.RouteName())
		ctx.Next()
	}, func(ctx iris.Context) {
		ctx.Writef("uid (uint): %d", ctx.Params().GetUintDefault("uid", 0))
	})

	app.Get("/u/{firstname:alphabetical}", func(ctx iris.Context) {
		ctx.Writef("before firstname (alphabetical), current route name: %s\n", ctx.RouteName())
		ctx.Next()
	}, func(ctx iris.Context) {
		ctx.Writef("firstname (alphabetical): %s", ctx.Params().Get("firstname"))
	})

	/*
		/u/some/path/here 对应 :path
		/u/abcd 对应 :alphabetical (如果不写 :alphabetical 默认是 :string)
		/u/42 对应 :uint (如果不写 :uint 默认是 :int)
		/u/-1 对应 :int (如果不写 :int 默认是 :string)
		/u/abcd123 对应 :string
	*/
	
	/*
		/u/some/path/here maps to :path
		/u/abcd maps to :alphabetical (if :alphabetical registered otherwise :string)
		/u/42 maps to :uint (if :uint registered otherwise :int)
		/u/-1 maps to :int (if :int registered otherwise :string)
		/u/abcd123 maps to :string
	*/

	// Pssst，别忘了使用动态路径示例获得更意外惊喜
	
	// Pssst, don't forget dynamic-path example for more "magic"!
	app.Get("/api/users/{userid:uint64 min(1)}", func(ctx iris.Context) {
		userID, err := ctx.Params().GetUint64("userid")
		if err != nil {
			ctx.Writef("error while trying to parse userid parameter," +
				"this will never happen if :uint64 is being used because if it's not a valid uint64 it will fire Not Found automatically.")
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		ctx.JSON(map[string]interface{}{
			//当然，您可以传递任何自定义的结构化go值。
			// you can pass any custom structured go value of course.
			"user_id": userID,
		})
	})
	// app.Post("/", func(ctx iris.Context){}) -> for POST http method.
	// app.Put("/", func(ctx iris.Context){})-> for "PUT" http method.
	// app.Delete("/", func(ctx iris.Context){})-> for "DELETE" http method.
	// app.Options("/", func(ctx iris.Context){})-> for "OPTIONS" http method.
	// app.Trace("/", func(ctx iris.Context){})-> for "TRACE" http method.
	// app.Head("/", func(ctx iris.Context){})-> for "HEAD" http method.
	// app.Connect("/", func(ctx iris.Context){})-> for "CONNECT" http method.
	// app.Patch("/", func(ctx iris.Context){})-> for "PATCH" http method.
	// app.Any("/", func(ctx iris.Context){}) for all http methods.

	//相同的路由可以对应多个不同的http 请求方法
	//您可以使用以下命令捕获任何路由创建错误：
	//路线， err := app.Get(...)
	//为路由设置名称：route.Name =“ myroute”
	
	// More than one route can contain the same path with a different http mapped method.
	// You can catch any route creation errors with:
	// route, err := app.Get(...)
	// set a name to a route: route.Name = "myroute"

	//您还可以按路径前缀对路由进行分组，共享中间件和完成的需要处理的动作。
	
	// You can also group routes by path prefix, sharing middleware(s) and done handlers.

	adminRoutes := app.Party("/admin", adminMiddleware)

	// Done 在ctx.Next()后面才会调用，所以下面 / 路由会调用一下ctx.Next()
	adminRoutes.Done(func(ctx iris.Context) { // executes always last if ctx.Next()
		ctx.Application().Logger().Infof("response sent to " + ctx.Path())
	})
	// adminRoutes.Layout("/views/layouts/admin.html")  // 为这些路由设置视图布局，请参view图示例。
	
	// adminRoutes.Layout("/views/layouts/admin.html") // set a view layout for these routes, see more at view examples.

	// GET: http://localhost:8080/admin
	adminRoutes.Get("/", func(ctx iris.Context) {
		// [...]
		ctx.StatusCode(iris.StatusOK) // default is 200 == iris.StatusOK
		ctx.HTML("<h1>Hello from admin/</h1>")

		// 为了执行路由组的 Done" Handler(s) 所以必须调用 ctx.Next() 
		ctx.Next() // in order to execute the party's "Done" Handler(s)
	})

	// GET: http://localhost:8080/admin/login
	adminRoutes.Get("/login", func(ctx iris.Context) {
		// [...]
	})
	// POST: http://localhost:8080/admin/login
	adminRoutes.Post("/login", func(ctx iris.Context) {
		// [...]
	})
	// 子域名比上面更容易, 执行要在host localhost或127.0.0.1
	// unix 路径在 /etc/hosts  windows 路径在 C:/windows/system32/drivers/etc/hosts
	
	// subdomains, easier than ever, should add localhost or 127.0.0.1 into your hosts file,
	// etc/hosts on unix or C:/windows/system32/drivers/etc/hosts on windows.

	//花括号是可选的，它只是样式的一种，以可视方式对路线进行分组。
	v1 := app.Party("v1.")
	{ // braces are optional, it's just type of style, to group the routes visually.

		// http://v1.localhost:8080
		
		//注意：对于版本特定的功能，请改用_examples /versioning。
		
		// Note: for versioning-specific features checkout the _examples/versioning instead.
		v1.Get("/", func(ctx iris.Context) {
			ctx.HTML(`Version 1 API. go to <a href="/api/users">/api/users</a>`)
		})

		usersAPI := v1.Party("/api/users")
		{
			// http://v1.localhost:8080/api/users
			usersAPI.Get("/", func(ctx iris.Context) {
				ctx.Writef("All users")
			})
			// http://v1.localhost:8080/api/users/42
			usersAPI.Get("/{userid:int}", func(ctx iris.Context) {
				ctx.Writef("user with id: %s", ctx.Params().Get("userid"))
			})
		}
	}
	// 通配符匹配子域名
	// wildcard subdomains.
	wildcardSubdomain := app.Party("*.")
	{
		wildcardSubdomain.Get("/", func(ctx iris.Context) {
			ctx.Writef("Subdomain can be anything, now you're here from: %s", ctx.Subdomain())
		})
	}

	return app
}

func main() {
	app := newApp()

	// http://localhost:8080
	// http://localhost:8080/home
	// http://localhost:8080/api/users/42
	// http://localhost:8080/admin
	// http://localhost:8080/admin/login
	//
	// http://localhost:8080/api/users/0
	// http://localhost:8080/api/users/blabla
	// http://localhost:8080/wontfound
	//
	// http://localhost:8080/u/abcd
	// http://localhost:8080/u/42
	// http://localhost:8080/u/-1
	// http://localhost:8080/u/abcd123
	// http://localhost:8080/u/some/path/here
	//
	// 修改host情况下:
	// if hosts edited:
	//  http://v1.localhost:8080
	//  http://v1.localhost:8080/api/users
	//  http://v1.localhost:8080/api/users/42
	//  http://anything.localhost:8080
	app.Run(iris.Addr(":8080"))
}

func adminMiddleware(ctx iris.Context) {
	// [...]
	//移至下一个处理程序，如果有任何身份验证逻辑，则不要调用这方法。
	ctx.Next() // to move to the next handler, or don't that if you have any auth logic.
}

func notFoundHandler(ctx iris.Context) {
	ctx.HTML("Custom route for 404 not found http code, here you can render a view, html, json <b>any valid response</b>.")
}

//注意：
//路由参数名称仅包含字母，符号则 _ ，数字将不被允许
//如果无法注册路由，则应用会在没有任何警告的情况下崩溃

//请参阅“file-server/single-page-application”以了解另一个功能“ WrapRouter”的工作方式。

// Notes:
// A path parameter name should contain only alphabetical letters, symbols, containing '_' and numbers are NOT allowed.
// If route failed to be registered, the app will panic without any warnings
// if you didn't catch the second return value(error) on .Handle/.Get....

// See "file-server/single-page-application" to see how another feature, "WrapRouter", works.
```
> `main_test.go`

```go
package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

// Shows a very basic usage of the httptest.
// The tests are written in a way to be easy to understand,
// for a more comprehensive testing examples check out the:
// _examples/routing/main_test.go,
// _examples/subdomains/www/main_test.go
// _examples/file-server and e.t.c.
// Almost every example which covers
// a new feature from you to learn
// contains a test file as well.
func TestRoutingBasic(t *testing.T) {
	expectedUResponse := func(paramName, paramType, paramValue string) string {
		s := fmt.Sprintf("before %s (%s), current route name: GET/u/{%s:%s}\n", paramName, paramType, paramName, paramType)
		s += fmt.Sprintf("%s (%s): %s", paramName, paramType, paramValue)
		return s
	}

	var (
		expectedNotFoundResponse = "Custom route for 404 not found http code, here you can render a view, html, json <b>any valid response</b>."

		expectedIndexResponse = "Hello from /"
		expectedHomeResponse  = `Same as app.Handle("GET", "/", [...])`

		expectedUpathResponse         = ":string, :int, :uint, :alphabetical and :path in the same path pattern."
		expectedUStringResponse       = expectedUResponse("username", "string", "abcd123")
		expectedUIntResponse          = expectedUResponse("id", "int", "-1")
		expectedUUintResponse         = expectedUResponse("uid", "uint", "42")
		expectedUAlphabeticalResponse = expectedUResponse("firstname", "alphabetical", "abcd")

		expectedAPIUsersIndexResponse = map[string]interface{}{"user_id": 42}

		expectedAdminIndexResponse = "<h1>Hello from admin/</h1>"

		expectedSubdomainV1IndexResponse                  = `Version 1 API. go to <a href="/api/users">/api/users</a>`
		expectedSubdomainV1APIUsersIndexResponse          = "All users"
		expectedSubdomainV1APIUsersIndexWithParamResponse = "user with id: 42"

		expectedSubdomainWildcardIndexResponse = "Subdomain can be anything, now you're here from: any-subdomain-here"
	)

	app := newApp()
	e := httptest.New(t, app)

	e.GET("/anotfound").Expect().Status(httptest.StatusNotFound).
		Body().Equal(expectedNotFoundResponse)

	e.GET("/").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedIndexResponse)
	e.GET("/home").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedHomeResponse)

	e.GET("/u/some/path/here").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedUpathResponse)
	e.GET("/u/abcd123").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedUStringResponse)
	e.GET("/u/-1").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedUIntResponse)
	e.GET("/u/42").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedUUintResponse)
	e.GET("/u/abcd").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedUAlphabeticalResponse)

	e.GET("/api/users/42").Expect().Status(httptest.StatusOK).
		JSON().Equal(expectedAPIUsersIndexResponse)

	e.GET("/admin").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedAdminIndexResponse)

	e.Request("GET", "/").WithURL("http://v1.example.com").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedSubdomainV1IndexResponse)

	e.Request("GET", "/api/users").WithURL("http://v1.example.com").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedSubdomainV1APIUsersIndexResponse)

	e.Request("GET", "/api/users/42").WithURL("http://v1.example.com").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedSubdomainV1APIUsersIndexWithParamResponse)

	e.Request("GET", "/").WithURL("http://any-subdomain-here.example.com").Expect().Status(httptest.StatusOK).
		Body().Equal(expectedSubdomainWildcardIndexResponse)
}
```