# go iris 视图 jet 模板第二示例
## 目录结构
> 主目录`template_jet_2`
```html
    —— views
        —— page.jet
    —— main.go
```
## 代码示例
> `views/page.jet`
```html
<a href="{{urlpath("my-page1")}}">/mypath</a>
<br />
<br/>
<a href="{{urlpath("my-page2","theParam1","theParam2")}}">/mypath2/{paramfirst}/{paramsecond}</a>
<br />
<br />

<a href="{{urlpath("my-page3", "theParam1", "theParam2AfterStatic")}}">/mypath3/{paramfirst}/statichere/{paramsecond}</a>
<br />
<br />

<a href="{{urlpath("my-page4", "theParam1", "theparam2AfterStatic",  "otherParam", "matchAnything")}}">
  /mypath4/{paramfirst}/statichere/{paramsecond}/{otherparam}/{something:path}</a>
<br />
<br />

<a href="{{urlpath("my-page5", "42", "theParam2Afterstatichere", "otherParam", "matchAnythingAfterStatic")}}">
  /mypath5/{paramfirst}/statichere/{paramsecond}/{otherparam}/anything/{anything:path}</a>
<br />
<br />

<a href={{urlpath("my-page6", .ParamsAsArray)}}>
  /mypath6/{paramfirst}/{paramsecond}/statichere/{paramThirdAfterStatic}
</a>
```
> `main.go`
```golang
//main包一个主要示例，说明如何命名您的路线并使用自定义'url path' Jet模板引擎
//
// Package main an example on how to naming your routes & use the custom 'url path' Jet Template Engine.
package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	app.RegisterView(iris.Jet("./views", ".jet").Reload(true))

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
	mypath5Route := app.Handle("GET", "/mypath5/{paramfirst:int}/statichere/{paramsecond}/{otherparam}/anything/{something:path}", writePathHandlerPage5)
	mypath5Route.Name = "my-page5"

	mypath6Route := app.Get("/mypath6/{paramfirst}/{paramsecond}/statichere/{paramThirdAfterStatic}", writePathHandler)
	mypath6Route.Name = "my-page6"

	app.Get("/", func(ctx iris.Context) {
		// for /mypath6...
		paramsAsArray := []string{"theParam1", "theParam2", "paramThirdAfterStatic"}
		ctx.ViewData("ParamsAsArray", paramsAsArray)
		if err := ctx.View("page.jet"); err != nil {
			panic(err)
		}
	})

	app.Get("/redirect/{namedRoute}", func(ctx iris.Context) {
		routeName := ctx.Params().Get("namedRoute")
		r := app.GetRoute(routeName)
		if r == nil {
			ctx.StatusCode(iris.StatusNotFound)
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

func writePathHandlerPage5(ctx iris.Context) {
	ctx.Writef("Hello from %s.\nparamfirst(int)=%d", ctx.Path(), ctx.Params().GetIntDefault("paramfirst", 0))
}
```