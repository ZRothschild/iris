# go iris 视图 html 模板第三示例
## 目录结构
> 主目录`template_html_3`
```html
    —— templates
        —— page.html
    —— main.go
```
## 代码示例
> `templates/page.html`
```html
<html>

<head>
  <title>template_html_3</title>
  <style>
    a {
      color: #0f7afc;
      border-bottom-color: rgba(15, 122, 252, 0.2);
      text-decoration: none
    }

    a:hover {
      color: #cf0000;
      border-bottom-color: rgba(208, 64, 0, 0.2);
      text-decoration: none
    }

    a:visited {
      color: #800080;
      border-bottom-color: rgba(128, 0, 128, 0.2);
      text-decoration: none
    }
  </style>
</head>

<body>

  <a href="{{urlpath "my-page1"}}">/mypath</a>
  <br />
  <br />

  <a href="{{urlpath "my-page2" "theParam1" "theParam2"}}">/mypath2/{paramfirst}/{paramsecond}</a>
  <br />
  <br />

  <a href="{{urlpath "my-page3" "theParam1" "theParam2AfterStatic"}}">/mypath3/{paramfirst}/statichere/{paramsecond}</a>
  <br />
  <br />

  <a href="{{urlpath "my-page4" "theParam1" "theparam2AfterStatic"  "otherParam"  "matchAnything"}}">
    /mypath4/{paramfirst}/statichere/{paramsecond}/{otherparam}/{something:path}</a>
  <br />
  <br />

  <a href="{{urlpath "my-page5" "theParam1" "theParam2Afterstatichere" "otherParam"  "matchAnythingAfterStatic"}}">
    /mypath5/{paramfirst}/statichere/{paramsecond}/{otherparam}/anything/{anything:path}</a>
  <br />
  <br />

  <a href={{urlpath "my-page6" .ParamsAsArray }}>
    /mypath6/{paramfirst}/{paramsecond}/statichere/{paramThirdAfterStatic}
  </a>
</body>

</html>
```
> `main.go`
```golang
//包一个主要示例，说明如何命名路由并使用自定义'url path'HTML模板引擎，其他模板引擎也是如此
// Package main an example on how to naming your routes & use the custom 'url path' HTML Template Engine, same for other template engines.
package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./templates", ".html").Reload(true))

	mypathRoute := app.Get("/mypath", writePathHandler)
	mypathRoute.Name = "my-page1"

	mypath2Route := app.Get("/mypath2/{paramfirst}/{paramsecond}", writePathHandler)
	mypath2Route.Name = "my-page2"

	mypath3Route := app.Get("/mypath3/{paramfirst}/statichere/{paramsecond}", writePathHandler)
	mypath3Route.Name = "my-page3"

	mypath4Route := app.Get("/mypath4/{paramfirst}/statichere/{paramsecond}/{otherparam}/{something:path}", writePathHandler)
	// same as: app.Get("/mypath4/:paramfirst/statichere/:paramsecond/:otherparam/*something", writePathHandler)
	mypath4Route.Name = "my-page4"

	// same with Handle/Func
	mypath5Route := app.Handle("GET", "/mypath5/{paramfirst}/statichere/{paramsecond}/{otherparam}/anything/{something:path}", writePathHandler)
	mypath5Route.Name = "my-page5"

	mypath6Route := app.Get("/mypath6/{paramfirst}/{paramsecond}/statichere/{paramThirdAfterStatic}", writePathHandler)
	mypath6Route.Name = "my-page6"

	app.Get("/", func(ctx iris.Context) {
		// for /mypath6...
		paramsAsArray := []string{"theParam1", "theParam2", "paramThirdAfterStatic"}
		ctx.ViewData("ParamsAsArray", paramsAsArray)
		if err := ctx.View("page.html"); err != nil {
			panic(err)
		}
	})

	app.Get("/redirect/{namedRoute}", func(ctx iris.Context) {
		routeName := ctx.Params().Get("namedRoute")
		r := app.GetRoute(routeName)
		if r == nil {
			ctx.StatusCode(404)
			ctx.Writef("Route with name %s not found", routeName)
			return
		}

		println("The path of " + routeName + "is: " + r.Path)
		//如果routeName == "my-page1"
		//打印：my-page1的路径是：/mypath
		//如果是使用命名参数的路径
		//然后使用"r.ResolvePath(paramValuesHere)"
		
		// if routeName == "my-page1"
		// prints: The path of of my-page1 is: /mypath
		// if it's a path which takes named parameters
		// then use "r.ResolvePath(paramValuesHere)"
		ctx.Redirect(r.Path)
		// http://localhost:8080/redirect/my-page1 will redirect to -> http://localhost:8080/mypath
	})

	// http://localhost:8080
	// http://localhost:8080/redirect/my-page1
	app.Run(iris.Addr(":8080"))
}

func writePathHandler(ctx iris.Context) {
	ctx.Writef("Hello from %s.", ctx.Path())
}
```