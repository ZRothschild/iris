# 域名重定向
## 目录结构
> 主目录`redirect`
```html
    —— hosts
    —— main.go
    —— main_test.go
```
## 代码示例
> `hosts`
```editorconfig
127.0.0.1	mydomain.com
127.0.0.1	www.mydomain.com

# Windows: Drive:/Windows/system32/drivers/etc/hosts, on Linux: /etc/hosts
```
> `main.go`
```go
//包main显示了如何使用app.WWW方法注册一个简单的'www'子域，
//该方法将注册一个路由器包装器，该包装器会将所有'mydomain.com'请求重定向到'www.mydomain.com'
// “hosts”文件，以了解如何在本地计算机上测试“ mydomain.com”。

// Package main shows how to register a simple 'www' subdomain,
// using the `app.WWW` method, which will register a router wrapper which will
// redirect all 'mydomain.com' requests to 'www.mydomain.com'.
// Check the 'hosts' file to see how to test the 'mydomain.com' on your local machine.
package main

import "github.com/kataras/iris/v12"

const addr = "mydomain.com:80"

func main() {
	app := newApp()

	// http(s)://mydomain.com,，将重定向到http(s)://www.mydomain.com
	//`www`变量是`app.Subdomain("www")`.

	// http(s)://mydomain.com, will be redirect to http(s)://www.mydomain.com.
	// The `www` variable is the `app.Subdomain("www")`.

	// app.WWW()包装路由器，以便它可以重定向所有传入的请求
	//来自'http(s)://mydomain.com/%path%'（缺少www）
	//到`http(s)://www.mydomain.com/%path%`.

	// app.WWW() wraps the router so it can redirect all incoming requests
	// that comes from 'http(s)://mydomain.com/%path%' (www is missing)
	// to `http(s)://www.mydomain.com/%path%`.

	// Try:
	// http://mydomain.com             -> http://www.mydomain.com
	// http://mydomain.com/users       -> http://www.mydomain.com/users
	// http://mydomain.com/users/login -> http://www.mydomain.com/users/login
	app.Run(iris.Addr(addr))
}

func newApp() *iris.Application {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("This will never be executed.")
	})

	www := app.Subdomain("www") // <- same as app.Party("www.") | 与app.Party("www.")相同
	www.Get("/", index)

	// www是一个`iris.Party`，请像使用路由一样对它进行分组。
	//与www.Party("/users").Get(...)相同

	// www is an `iris.Party`, use it like you already know, like grouping routes.

	www.PartyFunc("/users", func(p iris.Party) { // <- same as www.Party("/users").Get(...)
		p.Get("/", usersIndex)
		p.Get("/login", getLogin)
	})

	//将mydomain.com/%anypath%重定向到www.mydomain.com/%anypath%
	// 第一个参数是'from'，第二个参数是'to/target'

	// redirects mydomain.com/%anypath% to www.mydomain.com/%anypath%.
	// First argument is the 'from' and second is the 'to/target'.
	app.SubdomainRedirect(app, www)

	// SubdomainRedirect也适用于多级子域：

	// SubdomainRedirect works for multi-level subdomains as well:
	// subsub := www.Subdomain("subsub") // subsub.www.mydomain.com
	// subsub.Get("/", func(ctx iris.Context) { ctx.Writef("subdomain is: " + ctx.Subdomain()) })
	// app.SubdomainRedirect(subsub, www)

	//如果您需要将任何子域重定向到“ www”，则：
	// app.SubdomainRedirect（app.WildcardSubdomain（），www）
	//如果您需要从子域重定向到根域，则：
	// app.SubdomainRedirect（app.Subdomain（“ mysubdomain”），app）

	// If you need to redirect any subdomain to 'www' then:
	// app.SubdomainRedirect(app.WildcardSubdomain(), www)
	// If you need to redirect from a subdomain to the root domain then:
	// app.SubdomainRedirect(app.Subdomain("mysubdomain"), app)

	//注意，app.Party("mysubdomain.") 和app.Subdomain("mysubdomain")完全相同
	//不同之处在于第二个可以省略最后一个点('.')。
	
	// Note that app.Party("mysubdomain.") and app.Subdomain("mysubdomain")
	// is the same exactly thing, the difference is that the second can omit the last dot('.').

	return app
}

func index(ctx iris.Context) {
	ctx.Writef("This is the www.mydomain.com endpoint.")
}

func usersIndex(ctx iris.Context) {
	ctx.Writef("This is the www.mydomain.com/users endpoint.")
}

func getLogin(ctx iris.Context) {
	ctx.Writef("This is the www.mydomain.com/users/login endpoint.")
}
```
> `main_test.go`
```go
package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func TestSubdomainRedirectWWW(t *testing.T) {
	app := newApp()
	root := strings.TrimSuffix(addr, ":80")

	e := httptest.New(t, app)

	tests := []struct {
		path     string
		response string
	}{
		{"/", fmt.Sprintf("This is the www.%s endpoint.", root)},
		{"/users", fmt.Sprintf("This is the www.%s/users endpoint.", root)},
		{"/users/login", fmt.Sprintf("This is the www.%s/users/login endpoint.", root)},
	}

	for _, test := range tests {
		e.GET(test.path).Expect().Status(httptest.StatusOK).Body().Equal(test.response)
	}
}
```