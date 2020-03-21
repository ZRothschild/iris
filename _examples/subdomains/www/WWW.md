# IRIS www 域名使用
## 目录结构
> 主目录`WWW`
```html
    —— hosts
    —— main.go
    —— main_test.go
```
## 代码示例
> `main.go`
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
127.0.0.1		www.mydomain.com
#-END iris-
```
> `main.go`
```golang
package main

import (
	"github.com/kataras/iris/v12"
)

func newApp() *iris.Application {
	app := iris.New()

	app.Get("/", info)
	app.Get("/about", info)
	app.Get("/contact", info)

	app.PartyFunc("/api/users", func(r iris.Party) {
		r.Get("/", info)
		r.Get("/{id:uint64}", info)

		r.Post("/", info)

		r.Put("/{id:uint64}", info)
	})
	/* <- same as:
	 usersAPI := app.Party("/api/users")
	 {
		//这些括号只是语法糖,这种方法很少使用，
		//但是当您只想将范围变量用作该代码块时，可以使用它。

		// those brackets are just syntactic-sugar things.
		// This method is rarely used but you can make use of it when you want
	    // scoped variables to that code block only.
		usersAPI.Get/Post...
	 }
	 usersAPI.Get/Post...
	*/

	www := app.Party("www.")
	{
		//只是为了说明如何获取所有路由并将它们复制到另一方或子域：
		//获取到目前为止已注册的所有路由，包括所有"Parties"和子域：

		// Just to show how you can get all routes and copy them to another
		// party or subdomain:
		// Get all routes that are registered so far, including all "Parties" and subdomains:
		currentRoutes := app.GetRoutes()
		//同时将它们注册到 www subdomain/vhost：

		// Register them to the www subdomain/vhost as well:
		for _, r := range currentRoutes {
			www.Handle(r.Method, r.Tmpl().Src, r.Handlers...)
		}

		// http://www.mydomain.com/hi
		www.Get("/hi", func(ctx iris.Context) {
			ctx.Writef("hi from www.mydomain.com")
		})
	}
	//另请参见"subdomains/redirect"，以在子域之间注册重定向路由器包装，
	//即将mydomain.com更改为www.mydomain.com（像facebook一样由于SEO原因）。

	// See also the "subdomains/redirect" to register redirect router wrappers between subdomains,
	// i.e mydomain.com to www.mydomain.com (like facebook does for SEO reasons(;)).

	return app
}

func main() {
	app := newApp()
	// http://mydomain.com
	// http://mydomain.com/about
	// http://imydomain.com/contact
	// http://mydomain.com/api/users
	// http://mydomain.com/api/users/42

	// http://www.mydomain.com
	// http://www.mydomain.com/hi
	// http://www.mydomain.com/about
	// http://www.mydomain.com/contact
	// http://www.mydomain.com/api/users
	// http://www.mydomain.com/api/users/42
	if err := app.Run(iris.Addr("mydomain.com:80")); err != nil {
		panic(err)
	}
}

func info(ctx iris.Context) {
	method := ctx.Method()
	subdomain := ctx.Subdomain()
	path := ctx.Path()

	ctx.Writef("\nInfo\n\n")
	ctx.Writef("Method: %s\nSubdomain: %s\nPath: %s", method, subdomain, path)
}
```
> `main_test.go`
```golang
package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

type testRoute struct {
	path      string
	method    string
	subdomain string
}

func (r testRoute) response() string {
	msg := fmt.Sprintf("\nInfo\n\nMethod: %s\nSubdomain: %s\nPath: %s", r.method, r.subdomain, r.path)
	return msg
}

func TestSubdomainWWW(t *testing.T) {
	app := newApp()

	tests := []testRoute{
		// host
		{"/", "GET", ""},
		{"/about", "GET", ""},
		{"/contact", "GET", ""},
		{"/api/users", "GET", ""},
		{"/api/users/42", "GET", ""},
		{"/api/users", "POST", ""},
		{"/api/users/42", "PUT", ""},
		// www sub domain
		{"/", "GET", "www"},
		{"/about", "GET", "www"},
		{"/contact", "GET", "www"},
		{"/api/users", "GET", "www"},
		{"/api/users/42", "GET", "www"},
		{"/api/users", "POST", "www"},
		{"/api/users/42", "PUT", "www"},
	}

	host := "localhost:1111"
	e := httptest.New(t, app, httptest.Debug(false))

	for _, test := range tests {

		req := e.Request(test.method, test.path)
		if subdomain := test.subdomain; subdomain != "" {
			req.WithURL("http://" + subdomain + "." + host)
		}

		req.Expect().
			Status(httptest.StatusOK).
			Body().Equal(test.response())
	}
}
```