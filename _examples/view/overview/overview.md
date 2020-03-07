# go iris 视图使用回顾
## 目录结构
> 主目录`overview`
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
	<title>Hi iris</title>
</head>

<body>
	<h1>Hi {{.Name}} </h1>
</body>

</html>
```
> `main.go`
```golang
package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()

	//使用默认模板功能：
	//
	// with default template funcs:

	// - {{ urlpath "mynamedroute" "pathParameter_ifneeded" }}
	// - {{ render "header.html" }}
	// - {{ render_r "header.html" }} // partial relative path to current page
	// - {{ yield }}
	// - {{ current }}
	app.RegisterView(iris.HTML("./templates", ".html"))
	app.Get("/", func(ctx iris.Context) {
		// ..templates/hi.html中的.Name
		ctx.ViewData("Name", "iris") // the .Name inside the ./templates/hi.html
		//为大文件启用gzip
		ctx.Gzip(true)               // enable gzip for big files
		//使用相对于'./templates'的文件名渲染模板
		ctx.View("hi.html")          // render the template with the file name relative to the './templates'
	})

	// http://localhost:8080/
	app.Run(iris.Addr(":8080"))
}

/*
Note:
如果您想知道，视图引擎背后的代码来自"github.com/kataras/iris/v12/view"包，
也可以通过"github.com/kataras/iris/v12" 程序包访问引擎变量

In case you're wondering, the code behind the view engines derives from the "github.com/kataras/iris/v12/view" package,
access to the engines' variables can be granded by "github.com/kataras/iris/v12" package too.

    iris.HTML(...) is a shortcut of view.HTML(...)
    iris.Django(...)     >> >>      view.Django(...)
    iris.Pug(...)        >> >>      view.Pug(...)
    iris.Handlebars(...) >> >>      view.Handlebars(...)
    iris.Amber(...)      >> >>      view.Amber(...)
*/
```