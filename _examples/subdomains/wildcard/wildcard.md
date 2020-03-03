# 通配符动态匹配域名
## 目录结构
> 主目录`wildcard`
```html
    —— hosts
    —— main.go
```
## 代码示例
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
127.0.0.1		mydomain.com
127.0.0.1		username1.mydomain.com
127.0.0.1		username2.mydomain.com
127.0.0.1		username3.mydomain.com
127.0.0.1		username4.mydomain.com
127.0.0.1		username5.mydomain.com

#-END iris-
```
> `main.go`
```go
//包是一个关于如何捕获动态子域的通配符-通配符。在第一个示例（subdomains_1）中，我们看到了如何为静态子域（您知道将拥有的子域）创建路由
//在这里，我们将看到一个示例，说明如何捕获未知子域，动态子域，例如username.mydomain.com:8080。

// Package main an example on how to catch dynamic subdomains - wildcard.
// On the first example (subdomains_1) we saw how to create routes for static subdomains, subdomains you know that you will have.
// Here we will see an example how to catch unknown subdomains, dynamic subdomains, like username.mydomain.com:8080.
package main

import (
	"github.com/kataras/iris/v12"
)

//首先将动态通配符子域注册到服务器计算机(dns/...)，如果使用Windows，请检查./hosts
//运行此文件，然后尝试重定向：http://username1.mydomain.com:8080/，http://username2.mydomain.com:8080/，http://username1.mydomain.com/something，http：// username1.mydomain.com/something/sadsadsa

// register a dynamic-wildcard subdomain to your server machine(dns/...) first, check ./hosts if you use windows.
// run this file and try to redirect: http://username1.mydomain.com:8080/ , http://username2.mydomain.com:8080/ , http://username1.mydomain.com/something, http://username1.mydomain.com/something/sadsadsa

func main() {
	app := iris.New()

	/*
	请注意，您可以同时使用两种类型的子域（命名域和通配符(*.)
	admin.mydomain.com，以及其他方Party(*.)，但这不是此示例的目的

	Keep note that you can use both type of subdomains (named and wildcard(*.) )
	   admin.mydomain.com,  and for other the Party(*.) but this is not this example's purpose

	admin := app.Party("admin.")
	{
		// admin.mydomain.com
		admin.Get("/", func(ctx iris.Context) {
			ctx.Writef("INDEX FROM admin.mydomain.com")
		})
		// admin.mydomain.com/hey
		admin.Get("/hey", func(ctx iris.Context) {
			ctx.Writef("HEY FROM admin.mydomain.com/hey")
		})
		// admin.mydomain.com/hey2
		admin.Get("/hey2", func(ctx iris.Context) {
			ctx.Writef("HEY SECOND FROM admin.mydomain.com/hey")
		})
	}*/

	//没有顺序，您也可以在末尾注册子域。

	// no order, you can register subdomains at the end also.
	dynamicSubdomains := app.Party("*.")
	{
		dynamicSubdomains.Get("/", dynamicSubdomainHandler)

		dynamicSubdomains.Get("/something", dynamicSubdomainHandler)

		dynamicSubdomains.Get("/something/{paramfirst}", dynamicSubdomainHandlerWithParam)
	}

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Hello from mydomain.com path: %s", ctx.Path())
	})

	app.Get("/hello", func(ctx iris.Context) {
		ctx.Writef("Hello from mydomain.com path: %s", ctx.Path())
	})

	// http://mydomain.com:8080
	// http://username1.mydomain.com:8080
	// http://username2.mydomain.com:8080/something
	// http://username3.mydomain.com:8080/something/yourname

	//对于初学者：查看../hosts文件
	app.Run(iris.Addr("mydomain.com:8080")) // for beginners: look ../hosts file
}

func dynamicSubdomainHandler(ctx iris.Context) {
	username := ctx.Subdomain()
	ctx.Writef("Hello from dynamic subdomain path: %s, here you can handle the route for dynamic subdomains, handle the user: %s", ctx.Path(), username)
	//如果http://username4.mydomain.com:8080/打印：
	// 从动态子域路径：/中问好，在这里您可以处理动态子域的路由，请处理用户：username4
	
	// if  http://username4.mydomain.com:8080/ prints:
	// Hello from dynamic subdomain path: /, here you can handle the route for dynamic subdomains, handle the user: username4
}

func dynamicSubdomainHandlerWithParam(ctx iris.Context) {
	username := ctx.Subdomain()
	ctx.Writef("Hello from dynamic subdomain path: %s, here you can handle the route for dynamic subdomains, handle the user: %s", ctx.Path(), username)
	ctx.Writef("The paramfirst is: %s", ctx.Params().Get("paramfirst"))
}
```