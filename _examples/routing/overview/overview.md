# IRIS路由试用概览
## 目录结构
> 主目录`overview`
```html
    —— main.go
    —— public
        —— images
            —— favicon.ico
```
## 代码示例
> `main.go`
```golang
package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	// GET: http://localhost:8080
	app.Get("/", info)

	// GET: http://localhost:8080/profile/anyusername

	//是否要使用自定义正则表达式？
	//简单：app.Get("/profile/{username:string regexp(^[a-zA-Z ]+$)}")

	// Want to use a custom regex expression instead?
	// Easy: app.Get("/profile/{username:string regexp(^[a-zA-Z ]+$)}")
	app.Get("/profile/{username:string}", info)

	//如果缺少参数类型，则该字符串可以接受任何内容，
	//即：/{paramname} 与/{paramname:string}完全相同。
	//以下内容相同
	// {username:string}

	// If parameter type is missing then it's string which accepts anything,
	// i.e: /{paramname} it's exactly the same as /{paramname:string}.
	// The below is exactly the same as
	// {username:string}
	//
	// GET: http://localhost:8080/profile/anyusername/backups/any/number/of/paths/here
	app.Get("/profile/{username}/backups/{filepath:path}", info)

	// Favicon

	// GET: http://localhost:8080/favicon.ico
	app.Favicon("./public/images/favicon.ico")

	//静态文件
	//GET: http://localhost:8080/assets/css/bootstrap.min.css
	//映射到系统位置的./public/assets/css/bootstrap.min.css文件
	// GET: http://localhost:8080/assets/js/react.min.js
	//映射到系统位置的./public/assets/js/react.min.js文件

	// Static assets
	// GET: http://localhost:8080/assets/css/bootstrap.min.css
	//	    maps to ./public/assets/css/bootstrap.min.css file at system location.
	// GET: http://localhost:8080/assets/js/react.min.js
	//      maps to ./public/assets/js/react.min.js file at system location.
	app.HandleDir("/assets", "./public/assets")

	/* OR

	// GET: http://localhost:8080/js/react.min.js
	//映射到系统位置的./public/assets/js/react.min.js文件。
	app.HandleDir("/js", "./public/assets/js")

	// GET: http://localhost:8080/css/bootstrap.min.css
	//映射到系统位置的./public/assets/css/bootstrap.min.css文件。
	app.HandleDir("/css", "./public/assets/css")


	// GET: http://localhost:8080/js/react.min.js
	// 		maps to ./public/assets/js/react.min.js file at system location.
	app.HandleDir("/js", "./public/assets/js")

	// GET: http://localhost:8080/css/bootstrap.min.css
	// 		maps to ./public/assets/css/bootstrap.min.css file at system location.
	app.HandleDir("/css", "./public/assets/css")

	*/

	//分组
	// Grouping

	usersRoutes := app.Party("/users")
	// GET: http://localhost:8080/users/help
	usersRoutes.Get("/help", func(ctx iris.Context) {
		ctx.Writef("GET / -- fetch all users\n")
		ctx.Writef("GET /$ID -- fetch a user by id\n")
		ctx.Writef("POST / -- create new user\n")
		ctx.Writef("PUT /$ID -- update an existing user\n")
		ctx.Writef("DELETE /$ID -- delete an existing user\n")
	})

	// GET: http://localhost:8080/users
	usersRoutes.Get("/", func(ctx iris.Context) {
		ctx.Writef("get all users")
	})

	// GET: http://localhost:8080/users/42
	//iris版本7.0.5之后**/users/42和/users/help起作用**

	// **/users/42 and /users/help works after iris version 7.0.5**
	usersRoutes.Get("/{id:uint64}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetUint64("id")
		ctx.Writef("get user by id: %d", id)
	})

	// POST: http://localhost:8080/users
	usersRoutes.Post("/", func(ctx iris.Context) {
		username, password := ctx.PostValue("username"), ctx.PostValue("password")
		ctx.Writef("create user for username= %s and password= %s", username, password)
	})

	// PUT: http://localhost:8080/users
	usersRoutes.Put("/{id:uint64}", func(ctx iris.Context) {
		//或.Get获取其字符串表示形式。
		id, _ := ctx.Params().GetUint64("id") // or .Get to get its string represatantion.
		username := ctx.PostValue("username")
		ctx.Writef("update user for id= %d and new username= %s", id, username)
	})

	// DELETE: http://localhost:8080/users/42
	usersRoutes.Delete("/{id:uint64}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetUint64("id")
		ctx.Writef("delete user by id: %d", id)
	})

	//子域，取决于主机，如果使用它们，则必须编辑hosts文件或nginx/caddy的配置。
	//
	//在_examples/subdomains文件夹中查看更多子域示例。

	// Subdomains, depends on the host, you have to edit the hosts or nginx/caddy's configuration if you use them.
	//
	// See more subdomains examples at _examples/subdomains folder.
	adminRoutes := app.Party("admin.")

	// GET: http://admin.localhost:8080
	adminRoutes.Get("/", info)
	// GET: http://admin.localhost:8080/settings
	adminRoutes.Get("/settings", info)

	// Wildcard/dynamic subdomain
	dynamicSubdomainRoutes := app.Party("*.")

	// GET: http://any_thing_here.localhost:8080
	dynamicSubdomainRoutes.Get("/", info)

	app.Delete("/something", func(ctx iris.Context) {
		name := ctx.URLParam("name")
		ctx.Writef(name)
	})

	// GET: http://localhost:8080/
	// GET: http://localhost:8080/profile/anyusername
	// GET: http://localhost:8080/profile/anyusername/backups/any/number/of/paths/here

	// GET: http://localhost:8080/users/help
	// GET: http://localhost:8080/users
	// GET: http://localhost:8080/users/42
	// POST: http://localhost:8080/users
	// PUT: http://localhost:8080/users
	// DELETE: http://localhost:8080/users/42
	// DELETE: http://localhost:8080/something?name=iris

	// GET: http://admin.localhost:8080
	// GET: http://admin.localhost:8080/settings
	// GET: http://any_thing_here.localhost:8080
	app.Run(iris.Addr(":8080"))
}

func info(ctx iris.Context) {
	// http方法请求服务器的资源。
	method := ctx.Method()       // the http method requested a server's resource.
	//子域（如果有）。
	subdomain := ctx.Subdomain() // the subdomain, if any.
	//请求路径（没有协议和host）。
	// the request path (without scheme and host).
	path := ctx.Path()
	//如果不知道如何获取所有参数 字：
	// how to get all parameters, if we don't know
	// the names:
	paramsLen := ctx.Params().Len()

	ctx.Params().Visit(func(name string, value string) {
		ctx.Writef("%s = %s\n", name, value)
	})
	ctx.Writef("\nInfo\n\n")
	ctx.Writef("Method: %s\nSubdomain: %s\nPath: %s\nParameters length: %d", method, subdomain, path, paramsLen)
}
```