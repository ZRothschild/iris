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
