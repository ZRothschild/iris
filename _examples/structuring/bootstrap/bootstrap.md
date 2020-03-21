# IRIS `bootstrap`引导模式架构
## 目录结构
> 主目录`bootstrap.md`
```html
    —— bootstrap
        —— bootstrapper.go
    —— middleware
        —— identity
            —— identity.go
    —— public
        —— favicon.ico
    —— routes
        —— follower.go
        —— following.go
        —— index.go
        —— like.go
        —— routes.go
    —— views
        —— shared
            —— error.html
            —— layout.html
        —— index.html
    —— main.go
    —— main_test.go
```
## 目录结构图片
![目录结构图片](./folder_structure.png)
## 代码示例
> `bootstrap/bootstrapper.go`
```golang
package bootstrap

import (
	"time"

	"github.com/gorilla/securecookie"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/websocket"
)

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time

	Sessions *sessions.Sessions
}
// New返回一个新的Bootstrapper
// New returns a new Bootstrapper.
func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		Application:  iris.New(),
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}
// SetupViews加载模板
// SetupViews loads the templates.
func (b *Bootstrapper) SetupViews(viewsDir string) {
	b.RegisterView(iris.HTML(viewsDir, ".html").Layout("shared/layout.html"))
}
//SetupSessions可选地初始化会话
// SetupSessions initializes the sessions, optionally.
func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	b.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + b.AppName,
		Expires:  expires,
		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
	})
}
// SetupWebsockets设置websocket服务
// SetupWebsockets prepares the websocket server.
func (b *Bootstrapper) SetupWebsockets(endpoint string, handler websocket.ConnHandler) {
	ws := websocket.New(websocket.DefaultGorillaUpgrader, handler)

	b.Get(endpoint, websocket.Handler(ws))
}
// SetupErrorHandlers设置http错误处理程序
//`(context.StatusCodeNotSuccessful`，默认为<200 || > = 400，但可以更改）

// SetupErrorHandlers prepares the http error handlers
// `(context.StatusCodeNotSuccessful`,  which defaults to < 200 || >= 400 but you can change it).
func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		ctx.ViewData("Err", err)
		ctx.ViewData("Title", "Error")
		ctx.View("shared/error.html")
	})
}

const (
	// StaticAssets是公共文件目录（如图像，css, js）的根目录

	// StaticAssets is the root directory for public assets like images, css, js.
	StaticAssets = "./public/"
	// Favicon是我们应用程序 "StaticAssets" favicon路径

	// Favicon is the relative 9to the "StaticAssets") favicon path for our app.
	Favicon = "favicon.ico"
)
//Configure接受Configurator并在Bootstraper的上下文中运行它们

// Configure accepts configurations and runs them inside the Bootstraper's context.
func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}
// Bootstrap设置我们的应用程序
//
//返回自身

// Bootstrap prepares our application.
//
// Returns itself.
func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.SetupViews("./views")
	b.SetupSessions(24*time.Hour,
		[]byte("the-big-and-secret-fash-key-here"),
		[]byte("lot-secret-of-characters-big-too"),
	)
	b.SetupErrorHandlers()

	// static files 静态文件
	b.Favicon(StaticAssets + Favicon)
	b.HandleDir(StaticAssets[1:len(StaticAssets)-1], StaticAssets)
	//中间件，在静态文件之后

	// middleware, after static files
	b.Use(recover.New())
	b.Use(logger.New())

	return b
}
//Listen使用指定的“地址”启动http服务器

// Listen starts the http server with the specified "addr".
func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}
```
> `middleware/identity/identity.go`
```golang
package identity

import (
	"time"

	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/_examples/structuring/bootstrap/bootstrap"
)

// New返回一个新的处理程序，该处理程序添加了一些标题和视图数据
//描述应用程序信息比如所有者，启动时间

// New returns a new handler which adds some headers and view data
// describing the application, i.e the owner, the startup time.
func New(b *bootstrap.Bootstrapper) iris.Handler {
	return func(ctx iris.Context) {
		// response headers
		ctx.Header("App-Name", b.AppName)
		ctx.Header("App-Owner", b.AppOwner)
		ctx.Header("App-Since", time.Since(b.AppSpawnDate).String())

		ctx.Header("Server", "Iris: https://studyiris.com")
		//调用ctx.View或c.Tmpl = "$page.html"，渲染数据，然后调用ctx.Next() 继续往下执行

		// view data if ctx.View or c.Tmpl = "$page.html" will be called next.
		ctx.ViewData("AppName", b.AppName)
		ctx.ViewData("AppOwner", b.AppOwner)
		ctx.Next()
	}
}
//Configure创建一个新的身份中间件并将其注册到应用程序。

// Configure creates a new identity middleware and registers that to the app.
func Configure(b *bootstrap.Bootstrapper) {
	h := New(b)
	b.UseGlobal(h)
}
```
> `public/favicon.ico`
![项目图标](./public/favicon.ico)

> `routes/follower.go`
```golang
package routes

import (
	"github.com/kataras/iris/v12"
)

// GetFollowerHandler处理GET: /follower/{id}

// GetFollowerHandler handles the GET: /follower/{id}
func GetFollowerHandler(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")
	ctx.Writef("from "+ctx.GetCurrentRoute().Path()+" with ID: %d", id)
}
```
> `routes/following.go`
```golang
package routes

import (
	"github.com/kataras/iris/v12"
)

// GetFollowingHandler处理GET: /following/{id}

// GetFollowingHandler handles the GET: /following/{id}
func GetFollowingHandler(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")
	ctx.Writef("from "+ctx.GetCurrentRoute().Path()+" with ID: %d", id)
}
```
> `routes/index.go`
```golang
package routes

import (
	"github.com/kataras/iris/v12"
)

// GetIndexHandler处理GET: /

// GetIndexHandler handles the GET: /
func GetIndexHandler(ctx iris.Context) {
	ctx.ViewData("Title", "Index Page")
	ctx.View("index.html")
}
```
> `routes/like.go`
```golang
package routes

import (
	"github.com/kataras/iris/v12"
)

// GetLikeHandler处理GET: /like/{id}

// GetLikeHandler handles the GET: /like/{id}
func GetLikeHandler(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")
	ctx.Writef("from "+ctx.GetCurrentRoute().Path()+" with ID: %d", id)
}
```
> `routes/routes.go`
```golang
package routes

import (
	"github.com/kataras/iris/v12/_examples/structuring/bootstrap/bootstrap"
)
//Configure将必要的路由注册到应用程序

// Configure registers the necessary routes to the app.
func Configure(b *bootstrap.Bootstrapper) {
	b.Get("/", GetIndexHandler)
	b.Get("/follower/{id:long}", GetFollowerHandler)
	b.Get("/following/{id:long}", GetFollowingHandler)
	b.Get("/like/{id:long}", GetLikeHandler)
}
```
> `views/shared/error.html`
```html
<h1 class="text-danger">Error.</h1>
<h2 class="text-danger">An error occurred while processing your request.</h2>

<h3>{{.Err.status}}</h3>
<h4>{{.Err.message}}</h4>
```
> `views/shared/layout.html`
```html
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link rel="shortcut icon" type="image/x-icon" href="/favicon.ico" />
    <title>{{.Title}} - {{.AppName}}</title>

</head>

<body>
    <div>
        <!-- 在此处渲染当前模板 -->
        <!-- Render the current template here -->
        {{ yield }}
        <hr />
        <footer>
            <p>&copy; 2020 - {{.AppOwner}}</p>
        </footer>
    </div>
</body>

</html>
```
> `views/index.html`
```html
<h1>Welcome!!</h1>
```
> `main.go`
```golang
package main

import (
	"github.com/kataras/iris/v12/_examples/structuring/bootstrap/bootstrap"
	"github.com/kataras/iris/v12/_examples/structuring/bootstrap/middleware/identity"
	"github.com/kataras/iris/v12/_examples/structuring/bootstrap/routes"
)

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("Awesome App", "873908960@qq.com")
	app.Bootstrap()
	//identity.Configure中间件配置； routes.Configure路由配置 【进去他们的文件里面看看就懂了的】
	app.Configure(identity.Configure, routes.Configure)
	return app
}

func main() {
	app := newApp()
	app.Listen(":8080")
}
```
> `main_test.go`
```golang
package main

import (
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

// go test -v
func TestApp(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app.Application)

	// test our routes
	e.GET("/").Expect().Status(httptest.StatusOK)
	e.GET("/follower/42").Expect().Status(httptest.StatusOK).
		Body().Equal("from /follower/{id:long} with ID: 42")
	e.GET("/following/52").Expect().Status(httptest.StatusOK).
		Body().Equal("from /following/{id:long} with ID: 52")
	e.GET("/like/64").Expect().Status(httptest.StatusOK).
		Body().Equal("from /like/{id:long} with ID: 64")

	// test not found
	e.GET("/notfound").Expect().Status(httptest.StatusNotFound)
	expectedErr := map[string]interface{}{
		"app":     app.AppName,
		"status":  httptest.StatusNotFound,
		"message": "",
	}
	e.GET("/anotfoundwithjson").WithQuery("json", nil).
		Expect().Status(httptest.StatusNotFound).JSON().Equal(expectedErr)
}
```