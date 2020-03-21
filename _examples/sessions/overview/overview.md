# IRIS `seesion`会话使用简要介绍
## 目录结构
> 主目录`overview`
```html
    —— main.go
```
## 代码示例
> `main.go`
```golang
package main

import (
	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/sessions"
)

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID, AllowReclaim: true})
)

func secret(ctx iris.Context) {
	//检查用户是否已通过身份验证

	// Check if user is authenticated
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}
	//打印需要验证通过的信息

	// Print secret message
	ctx.WriteString("The cake is a lie!")
}

func login(ctx iris.Context) {
	session := sess.Start(ctx)
	//身份验证操作放在这里
	// Authentication goes here
	// ...
	//将用户设置为已认证
	// Set user as authenticated
	session.Set("authenticated", true)
}

func logout(ctx iris.Context) {
	session := sess.Start(ctx)
	//撤消用户身份已验证状态
	// Revoke users authentication
	session.Set("authenticated", false)
}

func main() {
	app := iris.New()

	app.Get("/secret", secret)
	app.Get("/login", login)
	app.Get("/logout", logout)

	app.Run(iris.Addr(":8080"))
}
```