# go iris 视图 html 模板引擎第4个示例
## 目录结构
> 主目录`template_html_4`
```html
    —— templates
        —— page.html
    —— hosts
    —— main.go
```
## 代码示例
> `templates/page.html`
```html
<!--  普通命名路由和动态子域命名路由之间的唯一区别是url的第一个参数
是子域部分，而不是命名参数 -->

<!-- the only difference between normal named routes and dynamic subdomains named routes is that the first argument of  url
is the subdomain part instead of named parameter-->

<a href="{{url "my-page1" "username1"}}">username1.127.0.0.1:8080/mypath</a>
<br />
<br />

<a href="{{url  "my-page2" "username2" "theParam1" "theParam2"}}">
    username2.127.0.0.1:8080/mypath2/{paramfirst}/{paramsecond}
</a>
<br />
<br />

<a href="{{url "my-page3" "username3" "theParam1" "theParam2AfterStatic"}}">
    username3.127.0.0.1:8080/mypath3/{paramfirst}/statichere/{paramsecond}
</a>
<br />
<br />

<a href="{{url "my-page4" "username4" "theParam1" "theparam2AfterStatic" "otherParam" "matchAnything"}}">
    username4.127.0.0.1:8080/mypath4/{paramfirst}/statichere/{paramsecond}/{otherParam}/{something:path}
</a>
<br />
<br />

<a href="{{url "my-page5" "username5" "theParam1" "theparam2AfterStatic" "otherParam" "matchAnything"}}">
    username5.127.0.0.1:8080/mypath5/{paramfirst}/statichere/{paramsecond}/{otherparam}/anything/{something:path}
</a>
<br/>
<br/>

<a href="{{url "my-page6" .ParamsAsArray }}">
    username5.127.0.0.1:8080/mypath6/{paramfirst}/{paramsecond}/staticParam/{paramThirdAfterStatic}
</a>
<br/>
<br/>

<a href="{{urlpath "my-page7" "theParam1" "theParam2" "theParam3" }}">
    mypath7/{paramfirst}/{paramsecond}/static/{paramthird}
</a>
<br/>
<br/>
```
> `hosts`
```editorconfig
# Copyright (c) 1993-2009 Microsoft Corp.
#
# This is a sample HOSTS file used by Microsoft TCP/IP for Windows.
#
# This file contains the mappings of IP addresses to host names. Each
# entry should be kept on an individual line. The IP address should
# be placed in the first column followed by the corresponding host name.
# The IP address and the host name should be separated by at least one
# space.
#
# Additionally, comments (such as these) may be inserted on individual
# lines or following the machine name denoted by a '#' symbol.
#
# For example:
#
#      102.54.94.97     rhino.acme.com          # source server
#       38.25.63.10     x.acme.com              # x client host

# localhost name resolution is handled within DNS itself.
127.0.0.1       localhost
::1             localhost
#-iris-For development machine, you have to configure your dns also for online, search google how to do it if you don't know

127.0.0.1		username1.127.0.0.1
127.0.0.1		username2.127.0.0.1
127.0.0.1		username3.127.0.0.1
127.0.0.1		username4.127.0.0.1
127.0.0.1		username5.127.0.0.1
# note that you can always use custom subdomains
#-END iris-

# Windows: Drive:/Windows/system32/drivers/etc/hosts, on Linux: /etc/hosts
```
> `main.go`
```golang
//包一个有关如何命名路由并使用自定义'url' HTML模板引擎（与其他模板引擎相同）的示例
// Package main an example on how to naming your routes & use the custom 'url' HTML Template Engine, same for other template engines.
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

const (
	host = "127.0.0.1:8080"
)

func main() {
	app := iris.New()
	//创建一个自定义的路由反向器，iris使您可以定义自己的主机和方案，
	//这在iris前面有nginx或caddy时很有用

	// create a custom path reverser, iris let you define your own host and scheme
	// which is useful when you have nginx or caddy in front of iris.
	rv := router.NewRoutePathReverser(app, router.WithHost(host), router.WithScheme("http"))
	//照常查找并定义我们的模板

	// locate and define our templates as usual.
	templates := iris.HTML("./templates", ".html")
	//添加一个自定义函数"url"，并将rv.URL作为其模板函数主体传递，
	// 因此{{url "routename" "paramsOrSubdomainAsFirstArgument"}}将在我们的模板中起作用

	// add a custom func of "url" and pass the rv.URL as its template function body,
	// so {{url "routename" "paramsOrSubdomainAsFirstArgument"}} will work inside our templates.
	templates.AddFunc("url", rv.URL)

	app.RegisterView(templates)

	//通配符子域，将捕获username1 .... username2 .... username3 ... username4 .... username5 ...
	//我们下面的链接是通过page.html的第一个参数（即子域）提供的

	// wildcard subdomain, will catch username1.... username2.... username3... username4.... username5...
	// that our below links are providing via page.html's first argument which is the subdomain.

	subdomain := app.Party("*.")

	mypathRoute := subdomain.Get("/mypath", emptyHandler)
	mypathRoute.Name = "my-page1"

	mypath2Route := subdomain.Get("/mypath2/{paramfirst}/{paramsecond}", emptyHandler)
	mypath2Route.Name = "my-page2"

	mypath3Route := subdomain.Get("/mypath3/{paramfirst}/statichere/{paramsecond}", emptyHandler)
	mypath3Route.Name = "my-page3"

	mypath4Route := subdomain.Get("/mypath4/{paramfirst}/statichere/{paramsecond}/{otherparam}/{something:path}", emptyHandler)
	mypath4Route.Name = "my-page4"

	mypath5Route := subdomain.Handle("GET", "/mypath5/{paramfirst}/statichere/{paramsecond}/{otherparam}/anything/{something:path}", emptyHandler)
	mypath5Route.Name = "my-page5"

	mypath6Route := subdomain.Get("/mypath6/{paramfirst}/{paramsecond}/staticParam/{paramThirdAfterStatic}", emptyHandler)
	mypath6Route.Name = "my-page6"

	app.Get("/", func(ctx iris.Context) {
		// for username5./mypath6...
		paramsAsArray := []string{"username5", "theParam1", "theParam2", "paramThirdAfterStatic"}
		ctx.ViewData("ParamsAsArray", paramsAsArray)
		if err := ctx.View("page.html"); err != nil {
			panic(err)
		}
	})
	//简单的路径，因此您可以在没有主机映射和子域的情况下对其进行测试，并在查看时使用{{urlpath ...}}，
	//以便向您展示如果您不希望整个方案和主机都可以使用它 成为网址的一部分

	// simple path so you can test it without host mapping and subdomains,
	// at view it make uses of {{urlpath ...}}
	// in order to showcase you that you can use it
	// if you don't want the entire scheme and host to be part of the url.
	app.Get("/mypath7/{paramfirst}/{paramsecond}/static/{paramthird}", emptyHandler).Name = "my-page7"

	// http://127.0.0.1:8080
	app.Run(iris.Addr(host))
}

func emptyHandler(ctx iris.Context) {
	ctx.Writef("Hello from subdomain: %s , you're in path:  %s", ctx.Subdomain(), ctx.Path())
}
// 注意：
//如果在{{url}}或{{urlpath}}上有一个空字符串，
//则表示args长度与路由的参数长度不匹配，
//或者传递的名称未找到该路由。

// Note:
// If you got an empty string on {{ url }} or {{ urlpath }} it means that
// args length are not aligned with the route's parameters length
// or the route didn't found by the passed name.
```