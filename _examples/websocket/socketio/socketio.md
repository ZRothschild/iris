# go iris websocket socketio 示例
## 目录结构
> 主目录`socketio`
```html
    —— asset
        —— index.html
    —— main.go
```
## 代码示例
> `asset/index.html`
```html
<!doctype html>
<html>

<head>
    <title>Socket.IO chat</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font: 13px Helvetica, Arial;
        }

        form {
            background: #000;
            padding: 3px;
            position: fixed;
            bottom: 0;
            width: 100%;
        }

        form input {
            border: 0;
            padding: 10px;
            width: 90%;
            margin-right: .5%;
        }

        form button {
            width: 9%;
            background: rgb(130, 224, 255);
            border: none;
            padding: 10px;
        }

        #messages {
            list-style-type: none;
            margin: 0;
            padding: 0;
        }

        #messages li {
            padding: 5px 10px;
        }

        #messages li:nth-child(odd) {
            background: #eee;
        }
    </style>
</head>

<body>
    <ul id="messages"></ul>
    <form action="">
        <input id="m" autocomplete="off" /><button>Send</button>
    </form>
    <script src="https://cdn.socket.io/socket.io-1.2.0.js"></script>
    <script src="https://code.jquery.com/jquery-1.11.1.js"></script>
    <script>
        var socket = io();
        // socket.emit('msg', 'hello');
        var s2 = io("/chat");
        socket.on('reply', function (msg) {
            $('#messages').append($('<li>').text(msg));
        });
        $('form').submit(function () {
            s2.emit('msg', $('#m').val(), function (data) {
                $('#messages').append($('<li>').text('ACK CALLBACK: ' + data));
            });
            socket.emit('notice', $('#m').val());
            $('#m').val('');
            return false;
        });
    </script>
</body>

</html>
```
> `main.go`
```golang
//包main运行一个基于go-socket.io的websocket服务器
//一个与Iris兼容的克隆：https://github.com/googollee/go-socket.io#example，
// 使用iris.FromStd来转换其处理程序
//
// Package main runs a go-socket.io based websocket server.
// An Iris compatible clone of: https://github.com/googollee/go-socket.io#example,
// use of `iris.FromStd` to convert its handler.
package main

import (
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	app.HandleMany("GET POST", "/socket.io/{any:path}", iris.FromStd(server))
	app.HandleDir("/", "./asset")
	app.Run(iris.Addr(":8000"),
		iris.WithoutPathCorrection,
		iris.WithoutServerError(iris.ErrServerClosed),
	)
}
```