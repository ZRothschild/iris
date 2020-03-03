# `go iris caddy`使用示例1
## 目录结构
> 主目录`server1`
```html
    —— Caddyfile
    —— views
        —— shared
            —— layout.html
         —— index.html
    —— main.go
```
## 代码示例
> `Caddyfile`
```editorconfig
example.com {
	header / Server "Iris"
	proxy / example.com:9091 # localhost:9091
}

api.example.com {
	header / Server "Iris"
	proxy / api.example.com:9092 # localhost:9092
}
```
> `views/shared/layout.html`
```html
<html>

<head>
    <title>{{.Layout.Title}}</title>
</head>

<body>
    {{ yield }}
</body>

</html>
```
> `views/index.html`
```html
<div>
    {{.Message}}
</div>
```
> `main.go`
```go
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := iris.New()

	templates := iris.HTML("./views", ".html").Layout("shared/layout.html")
	app.RegisterView(templates)

	mvc.New(app).Handle(new(Controller))

	// http://localhost:9091
	app.Run(iris.Addr(":9091"))
}
//布局包含shared/layout.html的所有绑定属性

// Layout contains all the binding properties for the shared/layout.html
type Layout struct {
	Title string
}
//Controller是我们的示例控制器，以请求为范围，每个请求都有其自己的实例

// Controller is our example controller, request-scoped, each request has its own instance.
type Controller struct {
	Layout Layout
}
// BeginRequest是客户端从此Controller的根路径请求时触发的第一个方法

// BeginRequest is the first method fired when client requests from this Controller's root path.
func (c *Controller) BeginRequest(ctx iris.Context) {
	c.Layout.Title = "Home Page"
}
// EndRequest是最后一个被触发的方法
//这里是为了完成BaseController
//以便告诉iris在main方法之前调用`BeginRequest`。

// EndRequest is the last method fired.
// It's here just to complete the BaseController
// in order to be tell iris to call the `BeginRequest` before the main method.
func (c *Controller) EndRequest(ctx iris.Context) {}

// Get handles GET http://localhost:9091
func (c *Controller) Get() mvc.View {
	return mvc.View{
		Name: "index.html",
		Data: iris.Map{
			"Layout":  c.Layout,
			"Message": "Welcome to my website!",
		},
	}
}
```