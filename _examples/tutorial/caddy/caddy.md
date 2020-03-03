# Caddy 与 Iris
`Caddyfile`展示了如何使用caddy来监听端口80和443，并坐在服务于其他端口（在这种情况下为9091和9092；请参见Caddyfile）的iris网络服务器前面
## 运行我们的两个Web服务器
1.转到`$GOPATH/src/github.com/kataras/iris/_examples/tutorial/caddy/server1`
2.打开一个终端窗口并执行`go run main.go`
3.转到`$GOPATH/src/github.com/kataras/iris/_examples/tutorial/caddy/server2`
4.打开一个新的终端窗口并执行`go run main.go`
## Caddy 安装
1.下载caddy: https://caddyserver.com/download
2.在这种情况下，将其内容提取到`Caddyfile`所在的位置，即 `$GOPATH/src/github.com/kataras/iris/_examples/tutorial/caddy`。
3.打开，阅读和修改`Caddyfile`，以自己了解配置服务器有多么容易
4.直接运行`caddy`或打开终端窗口并执行`caddy`
5.转到`https://example.com`和`https://api.example.com/user/42`
## Notes
Iris具有`app.Run(iris.AutoTLS(":443", "example.com", "mail@example.com"))` 
完全一样的东西，但是caddy是一个很棒的工具，当您在一台主机上运行多个Web服务器时，它可以为您提供帮助，例如Iirs，apache，tomcat。