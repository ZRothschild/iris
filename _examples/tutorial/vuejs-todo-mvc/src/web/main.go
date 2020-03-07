package main

import (
	"strings"

	"github.com/kataras/iris/v12/_examples/tutorial/vuejs-todo-mvc/src/todo"
	"github.com/kataras/iris/v12/_examples/tutorial/vuejs-todo-mvc/src/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/websocket"
)

func main() {
	app := iris.New()
	// 在public目录中提供我们的应用程序，公共文件夹包含客户端vue.js应用程序，
	// 这里不需要任何服务器端模板，
	// 实际上，如果您只使用vue而没有任何后端服务，
	// 则可以 在此行之后停止并启动服务器。

	// serve our app in public, public folder
	// contains the client-side vue.js application,
	// no need for any server-side template here,
	// actually if you're going to just use vue without any
	// back-end services, you can just stop afer this line and start the server.
	app.HandleDir("/", "./public")
	//配置http会话

	// configure the http sessions.
	sess := sessions.New(sessions.Config{
		Cookie: "iris_session",
	})
	//创建一个子路由器并注册http controllers

	// create a sub router and register the http controllers.
	todosRouter := app.Party("/todos")
	//创建针对/todos相对子路径的mvc应用程序

	// create our mvc application targeted to /todos relative sub path.
	todosApp := mvc.New(todosRouter)
	//这里的所有依赖项绑定...

	// any dependencies bindings here...
	todosApp.Register(
		todo.NewMemoryService(),
		sess.Start,
	)

	todosController := new(controllers.TodoController)
	//控制器注册在这里...

	// controllers registration here...
	todosApp.Handle(todosController)
	//为websocket控制器创建一个子mvc应用
	//继承父级的依赖项

	// Create a sub mvc app for websocket controller.
	// Inherit the parent's dependencies.
	todosWebsocketApp := todosApp.Party("/sync")
	todosWebsocketApp.HandleWebsocket(todosController).
		SetNamespace("todos").
		SetEventMatcher(func(methodName string) (string, bool) {
			return strings.ToLower(methodName), true
		})

	websocketServer := websocket.New(websocket.DefaultGorillaUpgrader, todosWebsocketApp)
	idGenerator := func(ctx iris.Context) string {
		id := sess.Start(ctx).ID()
		return id
	}
	todosWebsocketApp.Router.Get("/", websocket.Handler(websocketServer, idGenerator))
	//在http://localhost:8080启动Web服务器

	// start the web server at http://localhost:8080
	app.Run(iris.Addr(":8080"))
}
