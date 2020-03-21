# IRIS路由自定义上下文 新加字段
## 目录结构
> 主目录`new-implementation`
```html
    —— main.go
    —— main_test.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"sync"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)
//Owner是我们的应用程序结构，其中包含我们需要的方法或字段，
//将其视为*Context的所有者。

// Owner is our application structure, it contains the methods or fields we need,
// think it as the owner of our *Context.
type Owner struct {
	//在此定义全局字段
	//并共享给所有客户。

	// define here the fields that are global
	// and shared to all clients.
	sessionsManager *sessions.Sessions
}
//此包级变量“应用程序”将在上下文中用于与我们的全局应用程序进行通信。

// this package-level variable "application" will be used inside context to communicate with our global Application.
var owner = &Owner{
	sessionsManager: sessions.New(sessions.Config{Cookie: "mysessioncookie"}),
}
//上下文是我们的自定义上下文。
//让我们实现一个上下文，该上下文将使我们能够访问
//通过一个简单的`ctx.Session()`调用到客户端的Session。

// Context is our custom context.
// Let's implement a context which will give us access
// to the client's Session with a trivial `ctx.Session()` call.
type Context struct {
	iris.Context
	session *sessions.Session
}
//会话返回当前客户端的会话。
// Session returns the current client's session.
func (ctx *Context) Session() *sessions.Session {
	//如果我们在同一处理程序中多次调用`Session()`，这将对我们有帮助
	// this help us if we call `Session()` multiple times in the same handler
	if ctx.session == nil {
		//如果以前没有创建过，则开始新的会话。
		// start a new session if not created before.
		ctx.session = owner.sessionsManager.Start(ctx.Context)
	}

	return ctx.session
}
//粗体会向客户端发送粗体文本。
// Bold will send a bold text to the client.
func (ctx *Context) Bold(text string) {
	ctx.HTML("<b>" + text + "</b>")
}

var contextPool = sync.Pool{New: func() interface{} {
	return &Context{}
}}

func acquire(original iris.Context) *Context {
	ctx := contextPool.Get().(*Context)
	//将上下文设置为原始上下文，以便可以访问iris的实现。
	ctx.Context = original // set the context to the original one in order to have access to iris's implementation.
	//重置会话
	ctx.session = nil      // reset the session
	return ctx
}

func release(ctx *Context) {
	contextPool.Put(ctx)
}

//Handler会将我们的func(*Context)处理程序转换为Iris处理程序，
//用于HTTP API兼容。

// Handler will convert our handler of func(*Context) to an iris Handler,
// in order to be compatible with the HTTP API.
func Handler(h func(*Context)) iris.Handler {
	return func(original iris.Context) {
		ctx := acquire(original)
		h(ctx)
		release(ctx)
	}
}

func newApp() *iris.Application {
	app := iris.New()

	//像以前一样工作，唯一的不同
	//是原始context.Handler应该用我们的自定义包装
	//`Handler`函数。
	
	// Work as you did before, the only difference
	// is that the original context.Handler should be wrapped with our custom
	// `Handler` function.
	app.Get("/", Handler(func(ctx *Context) {
		ctx.Bold("Hello from our *Context")
	}))

	app.Post("/set", Handler(func(ctx *Context) {
		nameFieldValue := ctx.FormValue("name")
		ctx.Session().Set("name", nameFieldValue)
		ctx.Writef("set session = " + nameFieldValue)
	}))

	app.Get("/get", Handler(func(ctx *Context) {
		name := ctx.Session().GetString("name")
		ctx.Writef(name)
	}))

	return app
}

func main() {
	app := newApp()

	// GET: http://localhost:8080
	// POST: http://localhost:8080/set
	// GET: http://localhost:8080/get
	app.Run(iris.Addr(":8080"))
}
```

> `main_test.go`

```go
package main

import (
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func TestCustomContextNewImpl(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app, httptest.URL("http://localhost:8080"))

	e.GET("/").Expect().
		Status(httptest.StatusOK).
		ContentType("text/html").
		Body().Equal("<b>Hello from our *Context</b>")

	expectedName := "iris"
	e.POST("/set").WithFormField("name", expectedName).Expect().
		Status(httptest.StatusOK).
		Body().Equal("set session = " + expectedName)

	e.GET("/get").Expect().
		Status(httptest.StatusOK).
		Body().Equal(expectedName)
}
```