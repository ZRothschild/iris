package main

import (
	"errors"
	"time"

	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/badger"
)

func main() {
	db, err := badger.New("./data")
	if err != nil {
		panic(err)
	}
	//当按下control + C/cmd + C时关闭并解锁数据库

	// close and unlock the database when control+C/cmd+C pressed
	iris.RegisterOnInterrupt(func() {
		db.Close()
	})
	//如果应用程序出错，则关闭并解锁数据库
	defer db.Close() // close and unlock the database if application errored.

	sess := sessions.New(sessions.Config{
		Cookie:       "sessionscookieid",
		Expires:      45 * time.Minute, // <=0 means unlimited life. Defaults to 0. 小于等于0代表没有时间限制，默认是零
		AllowReclaim: true,
	})

	//
	// 重要：
	// IMPORTANT:
	//
	sess.UseDatabase(db)
	//其余代码保持不变

	// the rest of the code stays the same.
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("You should navigate to the /set, /get, /delete, /clear,/destroy instead")
	})
	app.Get("/set", func(ctx iris.Context) {
		s := sess.Start(ctx)
		//设置会话值

		// set session values
		s.Set("name", "iris")

		// test if set here
		ctx.Writef("All ok session value of the 'name' is: %s", s.GetString("name"))
	})

	app.Get("/set/{key}/{value}", func(ctx iris.Context) {
		key, value := ctx.Params().Get("key"), ctx.Params().Get("value")
		s := sess.Start(ctx)
		//设置会话值

		// set session values
		s.Set(key, value)
		//测试是否在这里设置

		// test if set here
		ctx.Writef("All ok session value of the '%s' is: %s", key, s.GetString(key))
	})

	app.Get("/get", func(ctx iris.Context) {
		//以字符串的形式获取特定键，如果找不到，则仅返回一个空字符串

		// get a specific key, as string, if no found returns just an empty string
		name := sess.Start(ctx).GetString("name")

		ctx.Writef("The 'name' on the /set was: %s", name)
	})

	app.Get("/get/{key}", func(ctx iris.Context) {
		//以字符串的形式获取特定键，如果找不到，则仅返回一个空字符串

		// get a specific key, as string, if no found returns just an empty string
		name := sess.Start(ctx).GetString(ctx.Params().Get("key"))

		ctx.Writef("The name on the /set was: %s", name)
	})

	app.Get("/delete", func(ctx iris.Context) {
		//删除特定的key

		// delete a specific key
		sess.Start(ctx).Delete("name")
	})

	app.Get("/clear", func(ctx iris.Context) {
		//删除所有key
		// removes all entries
		sess.Start(ctx).Clear()
	})

	app.Get("/destroy", func(ctx iris.Context) {
		//销毁，删除整个会话数据和cookie

		// destroy, removes the entire session data and cookie
		sess.Destroy(ctx)
	})

	app.Get("/update", func(ctx iris.Context) {
		//更新会根据会话的“Expires”字段重置过期时间。

		// updates resets the expiration based on the session's `Expires` field.
		if err := sess.ShiftExpiration(ctx); err != nil {
			if errors.Is(err, sessions.ErrNotFound) {
				ctx.StatusCode(iris.StatusNotFound)
			} else if errors.Is(err, sessions.ErrNotImplemented) {
				ctx.StatusCode(iris.StatusNotImplemented)
			} else {
				ctx.StatusCode(iris.StatusNotModified)
			}

			ctx.Writef("%v", err)
			ctx.Application().Logger().Error(err)
		}
	})

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
