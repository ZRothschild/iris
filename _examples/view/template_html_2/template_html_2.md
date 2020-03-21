# go iris 视图 html 模板第二示例
## 目录结构
> 主目录`template_html_2`
```html
    —— templates
        —— layouts
            —— layout.html
            —— mylayout.html
        —— partials
            —— page1_partial1.html
        —— page1.html
    —— main.go
```
## 代码示例
> `templates/layouts/layout.html`
```html
<html>
<head>
<title>Layout</title>

</head>
<body>
	<h1>This is the global layout</h1>
	<br />
	<!-- Render the current template here -->
	{{ yield }}
</body>
</html>
```
> `templates/layouts/mylayout.html`
```html
<html>
<head>
<title>my Layout</title>

</head>
<body>
	<h1>This is the layout for the /my/ and /my/other routes only</h1>
	<br />
	<!-- Render the current template here -->
	{{ yield }}
</body>
</html>
```
> `templates/partials/page1_partial1.html`
```html
<div style="background-color: white; color: red">
	<h1>Page 1's Partial 1</h1>
</div>
```
> `templates/page1.html`
```html
<div style="background-color: black; color: blue">

	<h1>Page 1 {{ greet "iris developer"}}</h1>

	{{ render "partials/page1_partial1.html"}}

</div>
```
> `main.go`
```golang
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
	//设置路由组的布局，.Layout应该在任何Get或其他Handle聚会方法之前

	// set a layout for a party, .Layout should be BEFORE any Get or other Handle party's method
	my := app.Party("/my").Layout("layouts/mylayout.html")
	//下面这两个都将使用layouts/mylayout.html作为其布局
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
```