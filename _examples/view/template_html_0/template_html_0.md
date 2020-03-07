# go iris 视图 html 模板第0个示例
## 目录结构
> 主目录`template_html_0.md`
```html
    —— templates
        —— hi.html
    —— main.go
```
## 代码示例
> `templates/hi.html`
```html
<html>
<head>
<title>{{.Title}}</title>
</head>
<body>
	<h1>Hi {{.Name}} </h1>
</body>
</html>
```
> `main.go`
```golang
package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New() //默认为这些 | defaults to these

	tmpl := iris.HTML("./templates", ".html")
	//根据每个请求重新加载模板（开发模式）
	tmpl.Reload(true) // reload templates on each request (development mode)
	//默认模板功能为：

	// default template funcs are:

	// - {{ urlpath "mynamedroute" "pathParameter_ifneeded" }}
	// - {{ render "header.html" }}
	// - {{ render_r "header.html" }} // partial relative path to current page
	// - {{ yield }}
	// - {{ current }}
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings " + s + "!"
	})
	app.RegisterView(tmpl)

	app.Get("/", hi)

	// http://localhost:8080
	//默认为该值，但您可以更改它
	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8")) // defaults to that but you can change it.
}

func hi(ctx iris.Context) {
	ctx.ViewData("Title", "Hi Page")
	ctx.ViewData("Name", "iris") // {{.Name}} will render: iris
	// ctx.ViewData("", myCcustomStruct{})
	ctx.View("hi.html")
}

/*
Note:

如果您想知道，视图引擎背后的代码来自"github.com/kataras/iris/v12/view"包，
也可以通过"github.com/kataras/iris/v12"程序包访问引擎变量

In case you're wondering, the code behind the view engines derives from the "github.com/kataras/iris/v12/view" package,
access to the engines' variables can be granded by "github.com/kataras/iris/v12" package too.

    iris.HTML(...) is a shortcut of view.HTML(...)
    iris.Django(...)     >> >>      view.Django(...)
    iris.Pug(...)        >> >>      view.Pug(...)
    iris.Handlebars(...) >> >>      view.Handlebars(...)
    iris.Amber(...)      >> >>      view.Amber(...)
*/
```