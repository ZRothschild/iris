# `Iris`动态记录线上访客量
## 目录结构
> 主目录`online-visitors`
```html
    —— static
        —— assets
            —— js
                —— visitors.js
    —— templates
        —— index.html
        —— other.html
    —— main.go
```
## 代码示例
> `static/assets/js/visitors.js`
```js
(function () {
  var events = {
    default: {
      _OnNamespaceConnected: function (ns, msg) {
        ns.joinRoom(PAGE_SOURCE);
      },
      _OnNamespaceDisconnect: function (ns, msg) {
        document.getElementById("online_views").innerHTML = "you've been disconnected";
      },
      onNewVisit: function (ns, msg) {
        var text = "1 online view";
        var onlineViews = Number(msg.Body);
        if (onlineViews > 1) {
          text = onlineViews + " online views";
        }
        document.getElementById("online_views").innerHTML = text;
      }
    }
  };

  neffos.dial("ws://localhost:8080/my_endpoint", events).then(function (client) {
    client.connect("default");
  });
})();
```
> `templates/index.html`
```html
<html>
<head>
    <title>Online visitors example</title>
    <style>
        body {
            margin: 0;
            font-family: -apple-system, "San Francisco", "Helvetica Neue", "Noto", "Roboto", "Calibri Light", sans-serif;
            color: #212121;
            font-size: 1.0em;
            line-height: 1.6;
        }
        .container {
            max-width: 750px;
            margin: auto;
            padding: 15px;
        }
        #online_views {
            font-weight: bold;
            font-size: 18px;
        }
    </style>
</head>
<body>
    <div class="container">
        <span id="online_views">1 online view</span>
    </div>
    <script type="text/javascript">
        /* take the page source from our passed struct  on .Render */
        var PAGE_SOURCE = {{ .PageID }}
    </script>
    <script src="https://cdn.jsdelivr.net/npm/neffos.js@latest/dist/neffos.min.js"></script>
    <script src="/js/visitors.js"></script>
</body>
</html>
```
> `templates/other.html`
```html
<html>
<head>
    <title>Different page, different results</title>
    <style>
        #online_views {
            font-weight: bold;
            font-size: 18px;
        }
    </style>
</head>
<body>
    <span id="online_views">1 online view</span>
    <script type="text/javascript">
        /* take the page source from our passed struct  on .Render */
        var PAGE_SOURCE = {{ .PageID }}
    </script>
    <script src="https://cdn.jsdelivr.net/npm/neffos.js@latest/dist/neffos.min.js"></script>
    <script src="/js/visitors.js"></script>
</body>

</html>
```
> `main.go` 结合 js 看看就明白了
```golang
package main

import (
	"fmt"
	"sync/atomic"

	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/websocket"
)

var events = websocket.Namespaces{
	"default": websocket.Events{
		websocket.OnRoomJoined: onRoomJoined,
		websocket.OnRoomLeft:   onRoomLeft,
	},
}

func main() {
	//初始化Web应用程序实例

	// init the web application instance
	// app := iris.New()
	app := iris.Default()
	//加载模板

	// load templates
	app.RegisterView(iris.HTML("./templates", ".html").Reload(true))
	//设置websocket服务器

	// setup the websocket server
	ws := websocket.New(websocket.DefaultGorillaUpgrader, events)

	app.Get("/my_endpoint", websocket.Handler(ws))

	//注册静态文件请求路径和系统目录

	// register static assets request path and system directory
	app.HandleDir("/js", "./static/assets/js")

	h := func(ctx iris.Context) {
		ctx.ViewData("", page{PageID: "index page"})
		ctx.View("index.html")
	}

	h2 := func(ctx iris.Context) {
		ctx.ViewData("", page{PageID: "other page"})
		ctx.View("other.html")
	}

	//打开一些浏览器标签页或窗口
	//并导航到
	// http://localhost:8080/ 和http://localhost:8080/other多次。
	//每个页面都有其自己的在线访客计数器。

	// Open some browser tabs/or windows
	// and navigate to
	// http://localhost:8080/ and http://localhost:8080/other multiple times.
	// Each page has its own online-visitors counter.
	app.Get("/", h)
	app.Get("/other", h2)
	app.Run(iris.Addr(":8080"))
}

type page struct {
	PageID string
}

type pageView struct {
	source string
	count  uint64
}

func (v *pageView) increment() {
	atomic.AddUint64(&v.count, 1)
}

func (v *pageView) decrement() {
	atomic.AddUint64(&v.count, ^uint64(0))
}

func (v *pageView) getCount() uint64 {
	return atomic.LoadUint64(&v.count)
}

type (
	pageViews []pageView
)

func (v *pageViews) Add(source string) {
	args := *v
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if kv.source == source {
			kv.increment()
			return
		}
	}

	c := cap(args)
	if c > n {
		args = args[:n+1]
		kv := &args[n]
		kv.source = source
		kv.count = 1
		*v = args
		return
	}

	kv := pageView{}
	kv.source = source
	kv.count = 1
	*v = append(args, kv)
}

func (v *pageViews) Get(source string) *pageView {
	args := *v
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if kv.source == source {
			return kv
		}
	}
	return nil
}

func (v *pageViews) Reset() {
	*v = (*v)[:0]
}

var v pageViews

func viewsCountBytes(viewsCount uint64) []byte {
	// *还有其他方法可以将uint64转换为[] byte

	// * there are other methods to convert uint64 to []byte
	return []byte(fmt.Sprintf("%d", viewsCount))
}

func onRoomJoined(ns *websocket.NSConn, msg websocket.Message) error {
	//这里的roomName是页面连接来源

	// the roomName here is the source.
	pageSource := string(msg.Room)

	v.Add(pageSource)

	viewsCount := v.Get(pageSource).getCount()
	if viewsCount == 0 {
		//这里的count应该总是> 0
		viewsCount++ // count should be always > 0 here
	}
	//在连接到该Room（来源的页面名称）的每个连接上触发"onNewVisit"客户端事件，
	// 并通知包括该连接在内的新访问（请参见第一个输入arg上的nil）

	// fire the "onNewVisit" client event
	// on each connection joined to this room (source page)
	// and notify of the new visit,
	// including this connection (see nil on first input arg).
	ns.Conn.Server().Broadcast(nil, websocket.Message{
		Namespace: msg.Namespace,
		Room:      pageSource,
		//触发"onNewVisit"客户端事件
		Event:     "onNewVisit", // fire the "onNewVisit" client event.
		Body:      viewsCountBytes(viewsCount),
	})

	return nil
}

func onRoomLeft(ns *websocket.NSConn, msg websocket.Message) error {
	//这里的roomName是来源的页面名称

	// the roomName here is the source.
	pageV := v.Get(msg.Room)
	if pageV == nil {
		//如果这个Room不是pageView源
		return nil // for any case that this room is not a pageView source
	}
	//递减-1此页面源的特定计数器

	// decrement -1 the specific counter for this page source.
	pageV.decrement()
	//在连接到该Room（源页面）的每个连接上触发"onNewVisit"客户端事件，并通知新的访问量（递减1）

	// fire the "onNewVisit" client event
	// on each connection joined to this room (source page)
	// and notify of the new, decremented by one, visits count.
	ns.Conn.Server().Broadcast(nil, websocket.Message{
		Namespace: msg.Namespace,
		Room:      msg.Room,
		Event:     "onNewVisit", //这个是visitors.js 触发的事件
		Body:      viewsCountBytes(pageV.getCount()),
	})

	return nil
}
```