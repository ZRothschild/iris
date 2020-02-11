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
	{ //花括号当然是可选的，这只是一种代码样式
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
	method := ctx.Method() // the http method requested a server's resource.
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
