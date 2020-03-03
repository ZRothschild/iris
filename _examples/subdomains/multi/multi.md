# 多域名使用【子域名】
## 目录结构
> 主目录`multi`
```html
    —— public
        —— assets
            —— images
                —— test.ico
        —— upload_resources
            —— favicon.ico
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

127.0.0.1 domain.local
127.0.0.1 system.domain.local
127.0.0.1 dashboard.domain.local

#-END iris-
```
> `main.go`
```go
package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	/*
	 * Setup static files 设置静态文件
	 */

	app.HandleDir("/assets", "./public/assets")
	app.HandleDir("/upload_resources", "./public/upload_resources")

	dashboard := app.Party("dashboard.")
	{
		dashboard.Get("/", func(ctx iris.Context) {
			ctx.Writef("HEY FROM dashboard")
		})
	}
	system := app.Party("system.")
	{
		system.Get("/", func(ctx iris.Context) {
			ctx.Writef("HEY FROM system")
		})
	}

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("HEY FROM frontend /")
	})
	// http://domain.local:80
	// http://dashboard.local
	// http://system.local

	//确保您在浏览器中添加"http"，因为.local是一个虚拟域，
	// 我们认为在这种情况下，您可以将任何语法正确的名称声明为iris的子域。

	// Make sure you prepend the "http" in your browser
	// because .local is a virtual domain we think to show case you
	// that you can declare any syntactical correct name as a subdomain in iris.
	
	//对于初学者：查看../hosts文件
	app.Run(iris.Addr("domain.local:80")) // for beginners: look ../hosts file
}
```