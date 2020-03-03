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
