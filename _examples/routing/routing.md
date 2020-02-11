# `route`基础用法
## 目录结构
> 主目录`routing`
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

/*
读：
"overview"
"basic"
"dynamic-path"
and "reverse"示例，如果你想能更好的使用iris

Read:
"overview"
"basic"
"dynamic-path"
and "reverse" examples if you want to release iris' real power.
*/

const maxBodySize = 1 << 20
const notFoundHTML = "<h1> custom http error page </h1>"

func registerErrors(app *iris.Application) {
	//设置自定义404处理程序

	// set a custom 404 handler
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.HTML(notFoundHTML)
	})
}

func registerGamesRoutes(app *iris.Application) {
	gamesMiddleware := func(ctx iris.Context) {
		ctx.Next()
	}

	// party只是一组具有相同前缀的路由
	//和中间件，即："/games"和gamesMiddleware。

	// party is just a group of routes with the same prefix
	// and middleware, i.e: "/games" and gamesMiddleware.
	games := app.Party("/games", gamesMiddleware)
	{ 	//花括号当然是可选的，这只是一种代码样式
		// braces are optional of course, it's just a style of code

		// "GET" method
		games.Get("/{gameID:uint64}/clans", h)
		games.Get("/{gameID:uint64}/clans/clan/{clanPublicID:uint64}", h)
		games.Get("/{gameID:uint64}/clans/search", h)

		// "PUT" method
		games.Put("/{gameID:uint64}/players/{clanPublicID:uint64}", h)
		games.Put("/{gameID:uint64}/clans/clan/{clanPublicID:uint64}", h)
		//切记："clanPublicID" 不应更改为具有相同前缀的其他路由。
		// remember: "clanPublicID" should not be changed to other routes with the same prefix.
		// "POST" method
		games.Post("/{gameID:uint64}/clans", h)
		games.Post("/{gameID:uint64}/players", h)
		games.Post("/{gameID:uint64}/clans/{clanPublicID:uint64}/leave", h)
		games.Post("/{gameID:uint64}/clans/{clanPublicID:uint64}/memberships/application", h)
		games.Post("/{gameID:uint64}/clans/{clanPublicID:uint64}/memberships/application/{action}", h) // {action} == {action:string}
		games.Post("/{gameID:uint64}/clans/{clanPublicID:uint64}/memberships/invitation", h)
		games.Post("/{gameID:uint64}/clans/{clanPublicID:uint64}/memberships/invitation/{action}", h)
		games.Post("/{gameID:uint64}/clans/{clanPublicID:uint64}/memberships/delete", h)
		games.Post("/{gameID:uint64}/clans/{clanPublicID:uint64}/memberships/promote", h)
		games.Post("/{gameID:uint64}/clans/{clanPublicID:uint64}/memberships/demote", h)

		gamesCh := games.Party("/challenge")
		{
			// games/challenge
			gamesCh.Get("/", h)

			gamesChBeginner := gamesCh.Party("/beginner")
			{
				// games/challenge/beginner/start
				gamesChBeginner.Get("/start", h)
				levelBeginner := gamesChBeginner.Party("/level")
				{
					// games/challenge/beginner/level/first
					levelBeginner.Get("/first", h)
				}
			}

			gamesChIntermediate := gamesCh.Party("/intermediate")
			{
				// games/challenge/intermediate
				gamesChIntermediate.Get("/", h)
				// games/challenge/intermediate/start
				gamesChIntermediate.Get("/start", h)
			}
		}

	}
}

func registerSubdomains(app *iris.Application) {
	mysubdomain := app.Party("mysubdomain.")
	// http://mysubdomain.myhost.com
	mysubdomain.Get("/", h)

	willdcardSubdomain := app.Party("*.")
	willdcardSubdomain.Get("/", h)
	willdcardSubdomain.Party("/party").Get("/", h)
}

func newApp() *iris.Application {
	app := iris.New()
	registerErrors(app)
	registerGamesRoutes(app)
	registerSubdomains(app)

	app.Handle("GET", "/healthcheck", h)

	//“ POST”方法
	//此处理程序从客户端/请求读取原始正文
	//并发送回相同的正文
	//记住，我们对那个正文有限制
	//保护自己免受“过热”的影响。

	// "POST" method
	// this handler reads raw body from the client/request
	// and sends back the same body
	// remember, we have limit to that body in order
	// to protect ourselves from "over heating".
	app.Post("/", iris.LimitRequestBodySize(maxBodySize), func(ctx iris.Context) {
		//获取请求主体

		// get request body
		b, err := ctx.GetBody()
		//如果较大，则发送错误的请求状态

		// if is larger then send a bad request status
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Writef(err.Error())
			return
		}
		// send back the post body
		ctx.Write(b)
	})

	app.HandleMany("POST PUT", "/postvalue", func(ctx iris.Context) {
		name := ctx.PostValueDefault("name", "iris")
		headervale := ctx.GetHeader("headername")
		ctx.Writef("Hello %s | %s", name, headervale)
	})

	return app
}

func h(ctx iris.Context) {
	// http方法请求服务器的资源
	method := ctx.Method()       // the http method requested a server's resource.
	//子域（如果有）
	subdomain := ctx.Subdomain() // the subdomain, if any.

	//请求路径（没有请求协议和主机）。

	// the request path (without scheme and host).
	path := ctx.Path()
	//如果不知道如何获取所有参数 名字：

	// how to get all parameters, if we don't know
	// the names:
	paramsLen := ctx.Params().Len()

	ctx.Params().Visit(func(name string, value string) {
		ctx.Writef("%s = %s\n", name, value)
	})
	ctx.Writef("Info\n\n")
	ctx.Writef("Method: %s\nSubdomain: %s\nPath: %s\nParameters length: %d", method, subdomain, path, paramsLen)
}

func main() {
	app := newApp()
	app.Logger().SetLevel("debug")

	/*
		// GET
		http://localhost:8080/healthcheck
		http://localhost:8080/games/42/clans
		http://localhost:8080/games/42/clans/clan/93
		http://localhost:8080/games/42/clans/search
		http://mysubdomain.localhost:8080/

		// PUT
		http://localhost:8080/postvalue
		http://localhost:8080/games/42/players/93
		http://localhost:8080/games/42/clans/clan/93

		// POST
		http://localhost:8080/
		http://localhost:8080/postvalue
		http://localhost:8080/games/42/clans
		http://localhost:8080/games/42/players
		http://localhost:8080/games/42/clans/93/leave
		http://localhost:8080/games/42/clans/93/memberships/application
		http://localhost:8080/games/42/clans/93/memberships/application/anystring
		http://localhost:8080/games/42/clans/93/memberships/invitation
		http://localhost:8080/games/42/clans/93/memberships/invitation/anystring
		http://localhost:8080/games/42/clans/93/memberships/delete
		http://localhost:8080/games/42/clans/93/memberships/promote
		http://localhost:8080/games/42/clans/93/memberships/demote

		//未找到触发
		// FIRE NOT FOUND
		http://localhost:8080/coudlntfound
	*/
	app.Run(iris.Addr(":8080"))
}
```

> `main_test.go`

```go
package main

import (
	"strconv"
	"strings"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func calculatePathAndResponse(method, subdomain, path string, paramKeyValue ...string) (string, string) {
	paramsLen := 0

	if l := len(paramKeyValue); l >= 2 {
		paramsLen = len(paramKeyValue) / 2
	}

	paramsInfo := ""
	if paramsLen > 0 {
		for i := 0; i < len(paramKeyValue); i++ {
			paramKey := paramKeyValue[i]
			i++
			if i >= len(paramKeyValue) {
				panic("paramKeyValue should be align with path parameters {} and must be placed in order")
			}

			paramValue := paramKeyValue[i]
			paramsInfo += paramKey + " = " + paramValue + "\n"

			beginParam := strings.IndexByte(path, '{')
			endParam := strings.IndexByte(path, '}')
			if beginParam == -1 || endParam == -1 {
				panic("something wrong with parameters, please define them in order")
			}

			path = path[:beginParam] + paramValue + path[endParam+1:]
		}
	}

	return path, paramsInfo + `Info

Method: ` + method + `
Subdomain: ` + subdomain + `
Path: ` + path + `
Parameters length: ` + strconv.Itoa(paramsLen)
}

type troute struct {
	method, subdomain, path string
	status                  int
	expectedBody            string
	contentType             string
}

func newTroute(method, subdomain, path string, status int, paramKeyValue ...string) troute {
	finalPath, expectedBody := calculatePathAndResponse(method, subdomain, path, paramKeyValue...)
	contentType := "text/plain; charset=UTF-8"

	if status == httptest.StatusNotFound {
		expectedBody = notFoundHTML
		contentType = "text/html; charset=UTF-8"
	}

	return troute{
		contentType:  contentType,
		method:       method,
		subdomain:    subdomain,
		path:         finalPath,
		status:       status,
		expectedBody: expectedBody,
	}
}

func TestRouting(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	tests := []troute{
		// GET
		newTroute("GET", "", "/healthcheck", httptest.StatusOK),
		newTroute("GET", "", "/games/{gameID}/clans", httptest.StatusOK, "gameID", "42"),
		newTroute("GET", "", "/games/{gameID}/clans/clan/{clanPublicID}", httptest.StatusOK, "gameID", "42", "clanPublicID", "93"),
		newTroute("GET", "", "/games/{gameID}/clans/search", httptest.StatusOK, "gameID", "42"),
		newTroute("GET", "", "/games/challenge", httptest.StatusOK),
		newTroute("GET", "", "/games/challenge/beginner/start", httptest.StatusOK),
		newTroute("GET", "", "/games/challenge/beginner/level/first", httptest.StatusOK),
		newTroute("GET", "", "/games/challenge/intermediate", httptest.StatusOK),
		newTroute("GET", "", "/games/challenge/intermediate/start", httptest.StatusOK),
		newTroute("GET", "mysubdomain", "/", httptest.StatusOK),
		newTroute("GET", "mywildcardsubdomain", "/", httptest.StatusOK),
		newTroute("GET", "mywildcardsubdomain", "/party", httptest.StatusOK),
		// PUT
		newTroute("PUT", "", "/games/{gameID}/players/{clanPublicID}", httptest.StatusOK, "gameID", "42", "clanPublicID", "93"),
		newTroute("PUT", "", "/games/{gameID}/clans/clan/{clanPublicID}", httptest.StatusOK, "gameID", "42", "clanPublicID", "93"),
		// POST
		newTroute("POST", "", "/games/{gameID}/clans", httptest.StatusOK, "gameID", "42"),
		newTroute("POST", "", "/games/{gameID}/players", httptest.StatusOK, "gameID", "42"),
		newTroute("POST", "", "/games/{gameID}/clans/{clanPublicID}/leave", httptest.StatusOK, "gameID", "42", "clanPublicID", "93"),
		newTroute("POST", "", "/games/{gameID}/clans/{clanPublicID}/memberships/application", httptest.StatusOK, "gameID", "42", "clanPublicID", "93"),
		newTroute("POST", "", "/games/{gameID}/clans/{clanPublicID}/memberships/application/{action}", httptest.StatusOK, "gameID", "42", "clanPublicID", "93", "action", "somethinghere"),
		newTroute("POST", "", "/games/{gameID}/clans/{clanPublicID}/memberships/invitation", httptest.StatusOK, "gameID", "42", "clanPublicID", "93"),
		newTroute("POST", "", "/games/{gameID}/clans/{clanPublicID}/memberships/invitation/{action}", httptest.StatusOK, "gameID", "42", "clanPublicID", "93", "action", "somethinghere"),
		newTroute("POST", "", "/games/{gameID}/clans/{clanPublicID}/memberships/delete", httptest.StatusOK, "gameID", "42", "clanPublicID", "93"),
		newTroute("POST", "", "/games/{gameID}/clans/{clanPublicID}/memberships/promote", httptest.StatusOK, "gameID", "42", "clanPublicID", "93"),
		newTroute("POST", "", "/games/{gameID}/clans/{clanPublicID}/memberships/demote", httptest.StatusOK, "gameID", "42", "clanPublicID", "93"),
		// POST: / will be tested alone
		// custom not found
		newTroute("GET", "", "/notfound", httptest.StatusNotFound),
		newTroute("POST", "", "/notfound2", httptest.StatusNotFound),
		newTroute("PUT", "", "/notfound3", httptest.StatusNotFound),
		newTroute("GET", "mysubdomain", "/notfound42", httptest.StatusNotFound),
	}

	for _, tt := range tests {
		et := e.Request(tt.method, tt.path)
		if tt.subdomain != "" {
			et.WithURL("http://" + tt.subdomain + ".localhost:8080")
		}
		et.Expect().Status(tt.status).Body().Equal(tt.expectedBody)
	}

	// test POST "/" limit data and post data return

	// test with small body
	e.POST("/").WithBytes([]byte("ok")).Expect().Status(httptest.StatusOK).Body().Equal("ok")
	// test with equal to max body size limit
	bsent := make([]byte, maxBodySize, maxBodySize)
	e.POST("/").WithBytes(bsent).Expect().Status(httptest.StatusOK).Body().Length().Equal(len(bsent))
	// test with larger body sent and wait for the custom response
	largerBSent := make([]byte, maxBodySize+1, maxBodySize+1)
	e.POST("/").WithBytes(largerBSent).Expect().Status(httptest.StatusBadRequest).Body().Equal("http: request body too large")

	// test the post value (both post and put) and headers.
	e.PUT("/postvalue").WithFormField("name", "test_put").
		WithHeader("headername", "headervalue_put").Expect().
		Status(httptest.StatusOK).Body().Equal("Hello test_put | headervalue_put")

	e.POST("/postvalue").WithFormField("name", "test_post").
		WithHeader("headername", "headervalue_post").Expect().
		Status(httptest.StatusOK).Body().Equal("Hello test_post | headervalue_post")
}
```