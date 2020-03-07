# go iris 上下文视图数据
## 目录结构
> 主目录`context-view-data`
```html
    —— templates
        —— layouts
            —— layout.html
        —— index.html
    —— main.go
```
## 代码示例
> `templates/layouts/layout.html`
```html
<html>
<head>
<title>My WebsiteLayout</title>

</head>
<body>
	<!-- Render the current template here -->
	{{ yield }}
</body>
</html>
```
> `templates/index.html`
```html
<h1>
	Title: {{.Title}}
</h1>
<h3>{{.BodyMessage}} </h3>

<hr/>

Current time: {{.CurrentTime}}
```
> `main.go`
```golang
package main

import (
	"time"

	"github.com/kataras/iris/v12"
)

const (
	DefaultTitle  = "My Awesome Site"
	DefaultLayout = "layouts/layout.html"
)

func main() {
	app := iris.New()
	//在os.Stdout上输出启动标语和错误日志

	// output startup banner and error logs on os.Stdout

	//将视图引擎目标设置为./templates文件夹

	// set the view engine target to ./templates folder
	app.RegisterView(iris.HTML("./templates", ".html").Reload(true))

	app.Use(func(ctx iris.Context) {
		//设置标题，当前时间和布局，以便在下一个处理程序调用.Render函数时使用

		// set the title, current time and a layout in order to be used if and when the next handler(s) calls the .Render function
		ctx.ViewData("Title", DefaultTitle)
		now := time.Now().Format(ctx.Application().ConfigurationReadOnly().GetTimeFormat())
		ctx.ViewData("CurrentTime", now)
		ctx.ViewLayout(DefaultLayout)

		ctx.Next()
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("BodyMessage", "a sample text here... set by the route handler")
		if err := ctx.View("index.html"); err != nil {
			ctx.Application().Logger().Infof(err.Error())
		}
	})

	app.Get("/about", func(ctx iris.Context) {
		ctx.ViewData("Title", "My About Page")
		ctx.ViewData("BodyMessage", "about text here... set by the route handler")

		//相同的文件，只是为了保持简单

		// same file, just to keep things simple.
		if err := ctx.View("index.html"); err != nil {
			ctx.Application().Logger().Infof(err.Error())
		}
	})

	// http://localhost:8080
	// http://localhost:8080/about
	app.Run(iris.Addr(":8080"))
}

//注意：ViewData("", myCustomStruct{})会将此myCustomStruct值设置为根绑定数据，
//因此任何View("other", "otherValue")都可能失败
//清除绑定数据：ctx.Set(ctx.Application().ConfigurationReadOnly().GetViewDataContextKey(), nil)

// Notes: ViewData("", myCustomStruct{}) will set this myCustomStruct value as a root binding data,
// so any View("other", "otherValue") will probably fail.
// To clear the binding data: ctx.Set(ctx.Application().ConfigurationReadOnly().GetViewDataContextKey(), nil)
```