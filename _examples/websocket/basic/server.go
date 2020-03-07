package main

import (
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	//当"enableJWT"常量为true时使用：

	// Used when "enableJWT" constant is true:
	"github.com/iris-contrib/middleware/jwt"
)

//值也应与客户端匹配

// values should match with the client sides as well.
const enableJWT = true
const namespace = "default"

//如果名称空间为空，则可以简单地使用websocket.Events {...}

// if namespace is empty then simply websocket.Events{...} can be used instead.
var serverEvents = websocket.Namespaces{
	namespace: websocket.Events{
		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			//通过`websocket.GetContext`可以获取iris的`Context`

			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			ctx := websocket.GetContext(nsConn.Conn)

			log.Printf("[%s] connected to namespace [%s] with IP [%s]",
				nsConn, msg.Namespace,
				ctx.RemoteAddr())
			return nil
		},
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
			return nil
		},
		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// room.String() 返回 -> NSConn.String() 返回 -> Conn.String() 返回 -> Conn.ID()

			// room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
			log.Printf("[%s] sent: %s", nsConn, string(msg.Body))
			// 使用以下命令将消息写回客户端消息所有者：
			// nsConn.Emit("chat", msg)
			// 使用以下命令向除此客户端之外的所有用户写消息：

			// Write message back to the client message owner with:
			// nsConn.Emit("chat", msg)
			// Write message to all except this client with:
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
	},
}

func main() {
	app := iris.New()
	websocketServer := websocket.New(
		/* 也可以使用DefaultGobwasUpgrader */
		websocket.DefaultGorillaUpgrader, /* DefaultGobwasUpgrader can be used too. */
		serverEvents)

	j := jwt.New(jwt.Config{
		//通过“token” URL提取，
		//因此客户端应使用ws://localhost:8080/echo?token=$token请求

		// Extract by the "token" url,
		// so the client should dial with ws://localhost:8080/echo?token=$token
		Extractor: jwt.FromParameter("token"),

		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// 设置后，中间件将验证令牌是否已使用特定的签名算法进行签名
		// 如果签名方法不是恒定的，则可使用`Config.ValidationKeyGetter`回调字段来实施其他检查，
		// 这对于避免此处所述的安全问题很重要：
		// https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/

		// When set, the middleware verifies that tokens are signed
		// with the specific signing algorithm
		// If the signing method is not constant the
		// `Config.ValidationKeyGetter` callback field can be used
		// to implement additional checks
		// Important to avoid security issues described here:
		// https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	idGen := func(ctx iris.Context) string {
		if username := ctx.GetHeader("X-Username"); username != "" {
			return username
		}

		return websocket.DefaultIDGenerator(ctx)
	}
	//通过可选的自定义ID生成器为ws://localhost:8080/echo的端点提供服务

	// serves the endpoint of ws://localhost:8080/echo
	// with optional custom ID generator.
	websocketRoute := app.Get("/echo", websocket.Handler(websocketServer, idGen))

	if enableJWT {
		//注册jwt中间件（握手时）：

		// Register the jwt middleware (on handshake):
		websocketRoute.Use(j.Serve)

		// 或|OR
		//
		//通过websocket连接或任何事件通过jwt中间件检查令牌：
		//
		// Check for token through the jwt middleware
		// on websocket connection or on any event:

		/* websocketServer.OnConnect = func(c *websocket.Conn) error {
		ctx := websocket.GetContext(c)
		if err := j.CheckJWT(ctx); err != nil {
			//将在客户端上发送上述错误，并且完全不允许它连接到websocket服务器

			// will send the above error on the client
			// and will not allow it to connect to the websocket server at all.
			return err
		}

		user := ctx.Values().Get("jwt").(*jwt.Token)
		// or just: user := j.Get(ctx)

		log.Printf("This is an authenticated request\n")
		log.Printf("Claim content:")
		log.Printf("%#+v\n", user.Claims)

		log.Printf("[%s] connected to the server", c.ID())

		return nil
		} */
	}
	//提供基于浏览器的websocket客户端

	// serves the browser-based websocket client.
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("./browser/index.html", false)
	})
	//提供npm浏览器websocket客户端用法示例

	// serves the npm browser websocket client usage example.
	app.HandleDir("/browserify", "./browserify")

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
