
# go iris 安全Websocket | Secure Websockets

1. 通过`https://` (`iris.Run(iris.TLS)` 或自定义的`iris.Listener(...)`运行服务器
2. 整个应用程序（包括websocket端）内部都没有变化
3. 客户端必须通过`wss://`前缀（而不是非安全的`wss://`），例如`wss://example.com/echo`，拨打Websocket服务器端点（即`/echo`）`/echo`
4. 准备出发

1. Run your server through `https://` (`iris.Run(iris.TLS)` or `iris.Run(iris.AutoTLS)` or a custom `iris.Listener(...)`)
2. Nothing changes inside the whole app, including the websocket side
3. The clients must dial the websocket server endpoint (i.e `/echo`) via `wss://` prefix (instead of the non-secure `ws://`), for example `wss://example.com/echo`
4. Ready to GO.
