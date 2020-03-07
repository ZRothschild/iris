package main

import (
	"github.com/kataras/iris/v12"
)

/*
您需要先构建hello.wasm，下载go1.11并执行以下命令：
$ cd client && GOARCH=wasm GOOS=js /home/$yourname/go1.11/bin/go build -o hello.wasm hello_go111.go

You need to build the hello.wasm first, download the go1.11 and execute the below command:
$ cd client && GOARCH=wasm GOOS=js /home/$yourname/go1.11/bin/go build -o hello.wasm hello_go111.go
*/

func main() {
	app := iris.New()
	//我们可以像例子一样，为您的资源提供服务，绝不在生产环境中包含.go文件

	// we could serve your assets like this the shake of the example,
	// never include the .go files there in production.
	app.HandleDir("/", "./client")

	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("./client/hello.html", false) // true 适用于gzip | true for gzip.
	})
	//访问http://localhost:8080
	//您应该获得这样的html输出：
	//您好，当前时间是：2018-07-09 05：54：12.564 +0000 UTC m = + 0.003900161

	// visit http://localhost:8080
	// you should get an html output like this:
	// Hello, the current time is: 2018-07-09 05:54:12.564 +0000 UTC m=+0.003900161
	app.Run(iris.Addr(":8080"))
}
