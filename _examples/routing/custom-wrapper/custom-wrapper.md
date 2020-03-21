# IRIS自定义路由包装器
## 目录结构
> 主目录`custom-wrapper`
```html 
    ——public
        —— css
            —— main.css
        —— app.js
        —— index.html
    —— main.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"net/http"
	"strings"

	"github.com/kataras/iris/v12"
)
//在此示例中，您只会看到.WrapRouter的一个用例。
//您可以使用.WrapRouter来添加自定义逻辑，无论何时router不应该
//执行以执行已注册路由的处理程序。

// In this example you'll just see one use case of .WrapRouter.
// You can use the .WrapRouter to add custom logic when or when not the router should
// be executed in order to execute the registered routes' handlers.
func newApp() *iris.Application {
	app := iris.New()

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.HTML("<b>Resource Not found</b>")
	})

	app.Get("/profile/{username}", func(ctx iris.Context) {
		ctx.Writef("Hello %s", ctx.Params().Get("username"))
	})

	app.HandleDir("/", "./public")

	myOtherHandler := func(ctx iris.Context) {
		ctx.Writef("inside a handler which is fired manually by our custom router wrapper")
	}

	//使用原生的net/http handler 包装router
	//如果url不包含任何"." (即: .css, .js...)
	//（取决于应用程序，您可能需要添加更多文件服务器异常），
	//然后处理程序将执行负责
	//已注册的routes(像 "/" and "/profile/{username}")
	//如果没有，则它将基于根"/"路径提供文件。

	// wrap the router with a native net/http handler.
	// if url does not contain any "." (i.e: .css, .js...)
	// (depends on the app , you may need to add more file-server exceptions),
	// then the handler will execute the router that is responsible for the
	// registered routes (look "/" and "/profile/{username}")
	// if not then it will serve the files based on the root "/" path.
	app.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/other") {
			//获取并释放上下文以便使用它来执行
			//我们的自定义处理程序
			//记住：我们使用net/http.Handler是因为这里我们在路由器底层之前的"低层"中。

			// acquire and release a context in order to use it to execute
			// our custom handler
			// remember: we use net/http.Handler because here we are in the "low-level", before the router itself.
			ctx := app.ContextPool.Acquire(w, r)
			myOtherHandler(ctx)
			app.ContextPool.Release(ctx)
			return
		}
		//否则继续照常提供routes
		router.ServeHTTP(w, r) // else continue serving routes as usual.
	})

	return app
}

func main() {
	app := newApp()

	// http://localhost:8080
	// http://localhost:8080/index.html
	// http://localhost:8080/app.js
	// http://localhost:8080/css/main.css
	// http://localhost:8080/profile/anyusername
	// http://localhost:8080/other/random
	app.Run(iris.Addr(":8080"))

	//注意：在此示例中，我们只看到了一个用例，
	//您可能希望使用.WrapRouter或.Downgrade来绕过Iris的默认路由器，即：
	//您也可以使用该方法来设置自定义代理。

	// Note: In this example we just saw one use case,
	// you may want to .WrapRouter or .Downgrade in order to bypass the iris' default router, i.e:
	// you can use that method to setup custom proxies too.
}
```
> `main_test.go`

```go
package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

type resource string

func (r resource) String() string {
	return string(r)
}

func (r resource) strip(strip string) string {
	s := r.String()
	return strings.TrimPrefix(s, strip)
}

func (r resource) loadFromBase(dir string) string {
	filename := r.String()

	if filename == "/" {
		filename = "/index.html"
	}

	fullpath := filepath.Join(dir, filename)

	b, err := ioutil.ReadFile(fullpath)
	if err != nil {
		panic(fullpath + " failed with error: " + err.Error())
	}

	return string(b)
}

var urls = []resource{
	"/",
	"/index.html",
	"/app.js",
	"/css/main.css",
}

func TestCustomWrapper(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	for _, u := range urls {
		url := u.String()
		contents := u.loadFromBase("./public")

		e.GET(url).Expect().
			Status(httptest.StatusOK).
			Body().Equal(contents)
	}

	e.GET("/other/something").Expect().Status(httptest.StatusOK)
}
```
> `public/index.html`

```html
<html>

<head>
    <title>Index Page</title>
</head>

<body>
    <h1> Hello from index.html </h1>


    <script src="/app.js">  </script>
</body>

</html>
```
> `public/app.js`

```js
window.alert("app.js loaded from \"/");
```
> `public/css/main.css`

```css
body {
    background-color: black;
}
```