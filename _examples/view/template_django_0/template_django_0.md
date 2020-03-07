# go iris 视图 django 模板
## 目录结构
> 主目录`template_django_0`
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
<title>{{title}}</title>
</head>
<body>
	<h1>Hi {{name|capfirst}} </h1>
	
	<h2>{{greet(name)}}</h2>

	<h3>Server started about {{serverStartTime|timesince}}. Refresh the page to see different result</h3>
</body>
</html>
```
> `main.go`
```golang
package main

import (
	"time"

	"github.com/kataras/iris/v12"

	// optionally, registers filters like `timesince`.
	_ "github.com/iris-contrib/pongo2-addons"
)

var startTime = time.Now()

func main() {
	app := iris.New()

	tmpl := iris.Django("./templates", ".html")
	//根据每个请求重新加载模板（开发模式）
	tmpl.Reload(true)                             // reload templates on each request (development mode)
	tmpl.AddFunc("greet", func(s string) string { // {{greet(name)}}
		return "Greetings " + s + "!"
	})

	// tmpl.RegisterFilter("myFilter", myFilter) // {{"simple input for filter"|myFilter}}
	app.RegisterView(tmpl)

	app.Get("/", hi)

	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}

func hi(ctx iris.Context) {
	// ctx.ViewData("title", "Hi Page")
	// ctx.ViewData("name", "iris")
	// ctx.ViewData("serverStartTime", startTime)
	// or if you set all view data in the same handler you can use the
	// iris.Map/pongo2.Context/map[string]interface{}, look below:

	ctx.View("hi.html", iris.Map{
		"title":           "Hi Page",
		"name":            "iris",
		"serverStartTime": startTime,
	})
}
```