# go iris 视图 html 模板第五示例
## 目录结构
> 主目录`template_html_5`
```html
    —— views
        —— user
            —— index.html
        —— about.html
        —— home.html
        —— layout.html
    —— main.go
```
## 代码示例
> `views/user/index.html`
```html
{{ define "user/index-head"}}
<style type="text/css">
    body {
        background: red;
    }
</style>
{{ end }}
<div>
    Hello user index page
</div>
```
> `views/about.html`
```html
{{ define "about-head"}}
    <title>about page</title>
    <style type="text/css">
        body {
            background: #666;
        }
    </style>
{{ end }}

{{ define "about-body"}}
    extend body content in layout.
{{ end }}
<div>
    Hello about page
</div>
```
> `views/home.html`
```html
{{ define "home-head"}}
	<title>{{.title}}</title>
	<style type="text/css">
		body {
			background: #999;
		}
	</style>
{{ end }}
<div>
	Hello home page
</div>
```
> `views/layout.html`
```html
{{ define "home-head"}}
	<title>{{.title}}</title>
	<style type="text/css">
		body {
			background: #999;
		}
	</style>
{{ end }}
<div>
	Hello home page
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

	app.RegisterView(iris.HTML("./views", ".html").Layout("layout.html"))
	//提示：追加.Reload(true)以在每个请求上重新加载模板

	// TIP: append .Reload(true) to reload the templates on each request.

	app.Get("/home", func(ctx iris.Context) {
		ctx.ViewData("title", "Home page")
		ctx.View("home.html")
		
		//注意：您可以传递"layout" : "otherLayout.html"绕过配置的Layout属性或view.NoLayout禁用此渲染操作的布局。 第三个是可选参数

		// Note that: you can pass "layout" : "otherLayout.html" to bypass the config's Layout property
		// or view.NoLayout to disable layout on this render action.
		// third is an optional parameter
	})

	app.Get("/about", func(ctx iris.Context) {
		ctx.View("about.html")
	})

	app.Get("/user/index", func(ctx iris.Context) {
		ctx.View("user/index.html")
	})

	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}
```