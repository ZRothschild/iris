package main

import (
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

type clientPage struct {
	Title string
	Host  string
}

func main() {
	app := iris.New()
	//选择要提供模板的html引擎
	app.RegisterView(iris.HTML("./templates", ".html")) // select the html engine to serve templates

	// neffos的几乎所有功能都被禁用，因为当应用程序希望仅接受并发送原始websocket本地消息时，自定义消息无法通过
	// 仅允许本机消息是事实吗？ 当注册的名称空间只有一个且为空时，仅包含一个注册的事件，即`OnNativeMessage`
	// 如果使用`Events {...}`而不是`Namespaces{ "namespaceName": Events{...}}`，则命名空间为空的 ""

	// Almost all features of neffos are disabled because no custom message can pass
	// when app expects to accept and send only raw websocket native messages.
	// When only allow native messages is a fact?
	// When the registered namespace is just one and it's empty
	// and contains only one registered event which is the `OnNativeMessage`.
	// When `Events{...}` is used instead of `Namespaces{ "namespaceName": Events{...}}`
	// then the namespace is empty "".
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
	})

	ws.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to server!", c.ID())
		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] Disconnected from server", c.ID())
	}

	//提供我们的自定义javascript代码
	app.HandleDir("/js", "./static/js") // serve our custom javascript code.
	//在端点上注册服务器
	//请参阅websockets.html中的内联javascript代码，此端点用于连接到服务器

	// register the server on an endpoint.
	// see the inline javascript code i the websockets.html, this endpoint is used to connect to the server.
	app.Get("/my_endpoint", websocket.Handler(ws))

	app.Get("/", func(ctx iris.Context) {
		ctx.View("client.html", clientPage{"Client Page", "localhost:8080"})
	})

	//将一些浏览器窗口/选项卡定位到http://localhost:8080并发送一些消息，
	//请参见static/js/chat.js，
	//请注意，客户端仅使用浏览器的本机WebSocket API，而不使用neffos

	// Target some browser windows/tabs to http://localhost:8080 and send some messages,
	// see the static/js/chat.js,
	// note that the client is using only the browser's native WebSocket API instead of the neffos one.
	app.Run(iris.Addr(":8080"))
}
