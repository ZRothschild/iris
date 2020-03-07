package controllers

import (
	"github.com/kataras/iris/v12/_examples/tutorial/vuejs-todo-mvc/src/todo"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/websocket"
)

// TodoController是我们的TODO应用程序的Web控制器

// TodoController is our TODO app's web controller.
type TodoController struct {
	Service todo.Service

	Session *sessions.Session

	NS *websocket.NSConn
}

//在服务器运行之前以及路由和依赖项绑定之前，调用一次BeforeActivation
// 您可以将自定义内容绑定到控制器，添加新方法，添加中间件，向结构或方法添加依赖项等等

// BeforeActivation called once before the server ran, and before
// the routes and dependencies binded.
// You can bind custom things to the controller, add new methods, add middleware,
// add dependencies to the struct or the method(s) and more.
func (c *TodoController) BeforeActivation(b mvc.BeforeActivation) {
	//可以绑定到控制器的函数输入参数
	// （如果有）或struct字段（如果有）：
	b.Dependencies().Add(func(ctx iris.Context) (items []todo.Item) {
		ctx.ReadJSON(&items)
		return
	})
}

// Get处理GET: /todos路线‘

// Get handles the GET: /todos route.
func (c *TodoController) Get() []todo.Item {
	return c.Service.Get(c.Session.ID())
}

// PostItemResponse所有todo items保存动作后将以json形式返回的响应数据

// PostItemResponse the response data that will be returned as json
// after a post save action of all todo items.
type PostItemResponse struct {
	Success bool `json:"success"`
}

var emptyResponse = PostItemResponse{Success: false}

// Post处理POST: /todos路由

// Post handles the POST: /todos route.
func (c *TodoController) Post(newItems []todo.Item) PostItemResponse {
	if err := c.Service.Save(c.Session.ID(), newItems); err != nil {
		return emptyResponse
	}

	return PostItemResponse{Success: true}
}

func (c *TodoController) Save(msg websocket.Message) error {
	id := c.Session.ID()
	c.NS.Conn.Server().Broadcast(nil, websocket.Message{
		Namespace: msg.Namespace,
		Event:     "saved",
		To:        id,
		Body:      websocket.Marshal(c.Service.Get(id)),
	})

	return nil
}
