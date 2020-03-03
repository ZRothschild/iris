// Package main展示了如何创建简单的URL Shortener。
//
//文章：https：//medium.com/@kataras/a-url-shortener-service-using-go-iris-and-bolt-4182f0b00ae7
//
// Package main shows how you can create a simple URL Shortener.
//
// Article: https://medium.com/@kataras/a-url-shortener-service-using-go-iris-and-bolt-4182f0b00ae7
//
// $ go get github.com/etcd-io/bbolt
// $ go get github.com/iris-contrib/go.uuid
// $ cd $GOPATH/src/github.com/kataras/iris/_examples/tutorial/url-shortener
// $ go build
// $ ./url-shortener
package main

import (
	"fmt"
	"html/template"

	"github.com/kataras/iris/v12"
)

func main() {
	//为数据库分配一个变量，以便稍后使用

	// assign a variable to the DB so we can use its features later.
	db := NewDB("shortener.db")
	//将该数据库传递给我们的应用程序，以便以后可以使用其他数据库测试整个应用程序。

	// Pass that db to our app, in order to be able to test the whole app with a different database later on.
	app := newApp(db)

	//当服务器关闭时释放"db"连接

	// release the "db" connection when server goes off.
	iris.RegisterOnInterrupt(db.Close)

	app.Run(iris.Addr(":8080"))
}

func newApp(db *DB) *iris.Application {
	app := iris.Default() // or app := iris.New()

	//创建我们的工厂，该工厂是对象创建的管理
	//在我们的Web应用程序和数据库之间

	// create our factory, which is the manager for the object creation.
	// between our web app and the db.
	factory := NewFactory(DefaultGenerator, db)

	//通过HTML std视图引擎为"./templates" 目录的“ * .html”文件提供服务

	// serve the "./templates" directory's "*.html" files with the HTML std view engine.
	tmpl := iris.HTML("./templates", ".html").Reload(true)
	//在此处注册任何模板功能
	//
	//看./templates/index.html#L16

	// register any template func(s) here.
	//
	// Look ./templates/index.html#L16
	tmpl.AddFunc("IsPositive", func(n int) bool {
		if n > 0 {
			return true
		}
		return false
	})

	app.RegisterView(tmpl)
	//提供静态文件（css）

	// Serve static files (css)
	app.HandleDir("/static", "./resources")

	indexHandler := func(ctx iris.Context) {
		ctx.ViewData("URL_COUNT", db.Len())
		ctx.View("index.html")
	}
	app.Get("/", indexHandler)

	//通过在http://localhost:8080/u/dsaoj41u321dsa上使用的键来查找并执行短网址

	// find and execute a short url by its key
	// used on http://localhost:8080/u/dsaoj41u321dsa
	execShortURL := func(ctx iris.Context, key string) {
		if key == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		value := db.Get(key)
		if value == "" {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.Writef("Short URL for key: '%s' not found", key)
			return
		}

		ctx.Redirect(value, iris.StatusTemporaryRedirect)
	}
	app.Get("/u/{shortkey}", func(ctx iris.Context) {
		execShortURL(ctx, ctx.Params().Get("shortkey"))
	})

	app.Get("/u/3861bc4d-ca57-4cbc-9fe4-9e0e2b50fff4", func(ctx iris.Context) {
		fmt.Printf("%s\n", "testsssss")
	})

	//app.Get("/u/{shortkey}", func(ctx iris.Context) {
	//	execShortURL(ctx, ctx.Params().Get("shortkey"))
	//})

	app.Post("/shorten", func(ctx iris.Context) {
		formValue := ctx.FormValue("url")
		if formValue == "" {
			ctx.ViewData("FORM_RESULT", "You need to a enter a URL")
			ctx.StatusCode(iris.StatusLengthRequired)
		} else {
			key, err := factory.Gen(formValue)
			if err != nil {
				ctx.ViewData("FORM_RESULT", "Invalid URL")
				ctx.StatusCode(iris.StatusBadRequest)
			} else {
				if err = db.Set(key, formValue); err != nil {
					ctx.ViewData("FORM_RESULT", "Internal error while saving the URL")
					app.Logger().Infof("while saving URL: " + err.Error())
					ctx.StatusCode(iris.StatusInternalServerError)
				} else {
					ctx.StatusCode(iris.StatusOK)
					shortenURL := "http://" + app.ConfigurationReadOnly().GetVHost() + "/u/" + key
					ctx.ViewData("FORM_RESULT",
						template.HTML("<pre><a target='_new' href='"+shortenURL+"'>"+shortenURL+" </a></pre>"))
				}
			}
		}
		//没有重定向，我们需要FORM_RESULT
		indexHandler(ctx) // no redirect, we need the FORM_RESULT.
	})

	app.Post("/clear_cache", func(ctx iris.Context) {
		db.Clear()
		ctx.Redirect("/")
	})

	return app
}
