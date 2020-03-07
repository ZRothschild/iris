package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	tmpl := iris.HTML("./templates", ".html")
	tmpl.Layout("layouts/layout.html")
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings " + s + "!"
	})

	// $ go get -u github.com/go-bindata/go-bindata/...
	// $ go-bindata ./templates/...
	// $ go build
	// $ ./embedding-templates-into-app

	//不使用html文件，您可以删除文件夹并运行示例

	// html files are not used, you can delete the folder and run the example.
	tmpl.Binary(Asset, AssetNames) // <-- 重要 | IMPORTANT

	app.RegisterView(tmpl)

	app.Get("/", func(ctx iris.Context) {
		if err := ctx.View("page1.html"); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Writef(err.Error())
		}
	})
	//删除特定路由的布局

	// remove the layout for a specific route
	app.Get("/nolayout", func(ctx iris.Context) {
		ctx.ViewLayout(iris.NoLayout)
		if err := ctx.View("page1.html"); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Writef(err.Error())
		}
	})
	//设置组路由的布局，.Layout应该在任何Get或其他Handle聚会方法之前

	// set a layout for a party, .Layout should be BEFORE any Get or other Handle party's method
	my := app.Party("/my").Layout("layouts/mylayout.html")
	//这两个路由都将使用layouts/mylayout.html作为其布局
	{ // both of these will use the layouts/mylayout.html as their layout.
		my.Get("/", func(ctx iris.Context) {
			ctx.View("page1.html")
		})
		my.Get("/other", func(ctx iris.Context) {
			ctx.View("page1.html")
		})
	}

	// http://localhost:8080
	// http://localhost:8080/nolayout
	// http://localhost:8080/my
	// http://localhost:8080/my/other
	app.Run(iris.Addr(":8080"))
}

//注意新的Gophers：
//如示例注释所示，使用`go build`代替`go run main.go`
//否则会出现编译错误，这是Go事情；
//因为在package main中有多个文件

// Note for new Gophers:
// `go build` is used instead of `go run main.go` as the example comments says
// otherwise you will get compile errors, this is a Go thing;
// because you have multiple files in the `package main`.
