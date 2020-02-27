# `seesion`会话中间件
## 目录结构
> 主目录`middleware`
```html
    —— main.go

```
## 代码示例
> `main.go`

```go
package main

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type businessModel struct {
	Name string
}

func main() {
	app := iris.New()
	sess := sessions.New(sessions.Config{
		// Cookie字符串，会话的客户端Cookie名称，例如：“ mysessionid”
		//
		//默认为“ irissessionid”

		// Cookie string, the session's client cookie name, for example: "mysessionid"
		//
		// Defaults to "irissessionid"
		Cookie: "mysessionid",
		//时间是 time.Duration 从创建cookie的时间开始，可以持续多久？
		// 0表示没有过期
		// -1表示浏览器关闭时到期
		//或设置一个值，例如2小时：

		// it's time.Duration, from the time cookie is created, how long it can be alive?
		// 0 means no expire.
		// -1 means expire when browser closes
		// or set a value, like 2 hours:
		Expires: time.Hour * 2,
		//如果您想使不同子域上的cookie无效
		//同一台主机，然后启用它即可
		//默认为false

		// if you want to invalid cookies on different subdomains
		// of the same host, then enable it.
		// Defaults to false.
		DisableSubdomainPersistence: false,
	})
	//会话现在在处理程序中始终为非零
	app.Use(sess.Handler()) // session is always non-nil inside handlers now.

	app.Get("/", func(ctx iris.Context) {
		session := sessions.Get(ctx) // same as sess.Start(ctx, cookieOptions...)
		if session.Len() == 0 {
			ctx.HTML(`no session values stored yet. Navigate to: <a href="/set">set page</a>`)
			return
		}

		ctx.HTML("<ul>")
		session.Visit(func(key string, value interface{}) {
			ctx.HTML("<li> %s = %v </li>", key, value)
		})

		ctx.HTML("</ul>")
	})
	//设置会话值

	// set session values.
	app.Get("/set", func(ctx iris.Context) {
		session := sessions.Get(ctx)
		session.Set("name", "iris")
		//测试是否在此处设置

		// test if set here.
		ctx.Writef("All ok session set to: %s", session.GetString("name"))
		// Set将按原样设置值，
		//如果是切片或map
		//您将可以在上对其进行更改。直接获取！
		//请注意，我不建议在会话中既不保存切片又不保存大数据
		//，但如果确实需要，请使用SetImmutable而不是Set。
		//过多使用`SetImmutable`，速度较慢。
		//阅读有关可变和不可变go类型的更多信息：https://stackoverflow.com/a/8021081

		// Set will set the value as-it-is,
		// if it's a slice or map
		// you will be able to change it on .Get directly!
		// Keep note that I don't recommend saving big data neither slices or maps on a session
		// but if you really need it then use the `SetImmutable` instead of `Set`.
		// Use `SetImmutable` consistently, it's slower.
		// Read more about muttable and immutable go types: https://stackoverflow.com/a/8021081
	})

	app.Get("/set/{key}/{value}", func(ctx iris.Context) {
		key, value := ctx.Params().Get("key"), ctx.Params().Get("value")

		session := sessions.Get(ctx)
		session.Set(key, value)
		//测试是否在此处设置

		// test if set here
		ctx.Writef("All ok session value of the '%s' is: %s", key, session.GetString(key))
	})

	app.Get("/get", func(ctx iris.Context) {
		//获取特定值，例如字符串，
		//如果未找到，则仅返回一个空字符串。

		// get a specific value, as string,
		// if not found then it returns just an empty string.
		name := sessions.Get(ctx).GetString("name")

		ctx.Writef("The name on the /set was: %s", name)
	})

	app.Get("/delete", func(ctx iris.Context) {
		//删除特定的key

		// delete a specific key
		sessions.Get(ctx).Delete("name")
	})

	app.Get("/clear", func(ctx iris.Context) {
		//删除所有key对应的value

		// removes all entries.
		sessions.Get(ctx).Clear()
	})

	app.Get("/update", func(ctx iris.Context) {
		//更新过期时间

		// updates expire date.
		sess.ShiftExpiration(ctx)
	})

	app.Get("/destroy", func(ctx iris.Context) {
		//销毁，删除整个会话数据和cookie
		// sess.Destroy(ctx)
		// 要么

		// destroy, removes the entire session data and cookie
		// sess.Destroy(ctx)
		// or
		sessions.Get(ctx).Destroy()
	})
	//有关销毁的注意事项：
	//
	//您也可以使用以下命令销毁处理程序外部的会话：
	// sess.DestroyByID
	// sess.DestroyAll

	// Note about Destroy:
	//
	// You can destroy a session outside of a handler too, using the:
	// sess.DestroyByID
	// sess.DestroyAll

	//切记：切片和map是根据设计可变的
	//`SetImmutable`确保它们将被存储和接收
	//是不可变的，因此您不能错误地直接更改它们。
	//
	//使用`SetImmutable`，它比`Set`慢。
	//阅读有关可变和不可变go类型的更多信息：https://stackoverflow.com/a/8021081

	// remember: slices and maps are muttable by-design
	// The `SetImmutable` makes sure that they will be stored and received
	// as immutable, so you can't change them directly by mistake.
	//
	// Use `SetImmutable` consistently, it's slower than `Set`.
	// Read more about muttable and immutable go types: https://stackoverflow.com/a/8021081
	app.Get("/set_immutable", func(ctx iris.Context) {
		business := []businessModel{{Name: "Edward"}, {Name: "value 2"}}
		session := sessions.Get(ctx)
		session.SetImmutable("businessEdit", business)
		businessGet := session.Get("businessEdit").([]businessModel)
		//尝试更改它，如果我们使用`Set`而不是`SetImmutable`
		//更改将影响会话值"businessEdit"的数组，但现在不会。

		// try to change it, if we used `Set` instead of `SetImmutable` this
		// change will affect the underline array of the session's value "businessEdit", but now it will not.
		businessGet[0].Name = "Gabriel"
	})

	app.Get("/get_immutable", func(ctx iris.Context) {
		valSlice := sessions.Get(ctx).Get("businessEdit")
		if valSlice == nil {
			ctx.HTML("please navigate to the <a href='/set_immutable'>/set_immutable</a> first")
			return
		}

		firstModel := valSlice.([]businessModel)[0]
		// businessGet [0] .Name最初等于Edward

		// businessGet[0].Name is equal to Edward initially
		if firstModel.Name != "Edward" {
			panic("Report this as a bug, immutable data cannot be changed from the caller without re-SetImmutable")
		}

		ctx.Writef("[]businessModel[0].Name remains: %s", firstModel.Name)
		//名称应保持为"Edward"
		
		// the name should remains "Edward"
	})

	app.Run(iris.Addr(":8080"))
}
```