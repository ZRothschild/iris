# `iris`服务端缓存
## 目录结构
> 主目录`simple`
```html
    —— main.go
```
## 代码示例 
> `main.go`
```golang
package main

import (
	"time"

	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/cache"
)

var markdownContents = []byte(`## Hello Markdown

This is a sample of Markdown contents

 

Features
--------

All features of Sundown are supported, including:

*   **Compatibility**. The Markdown v1.0.3 test suite passes with
    the --tidy option.  Without --tidy, the differences are
    mostly in whitespace and entity escaping, where blackfriday is
    more consistent and cleaner.

*   **Common extensions**, including table support, fenced code
    blocks, autolinks, strikethroughs, non-strict emphasis, etc.

*   **Safety**. Blackfriday is paranoid when parsing, making it safe
    to feed untrusted user input without fear of bad things
    happening. The test suite stress tests this and there are no
    known inputs that make it crash.  If you find one, please let me
    know and send me the input that does it.

    NOTE: "safety" in this context means *runtime safety only*. In order to
    protect yourself against JavaScript injection in untrusted content, see
    [this example](https://github.com/russross/blackfriday#sanitize-untrusted-content).

*   **Fast processing**. It is fast enough to render on-demand in
    most web applications without having to cache the output.

*   **Routine safety**. You can run multiple parsers in different
    goroutines without ill effect. There is no dependence on global
    shared state.

*   **Minimal dependencies**. Blackfriday only depends on standard
    library packages in Go. The source code is pretty
    self-contained, so it is easy to add to any project, including
    Google App Engine projects.

*   **Standards compliant**. Output successfully validates using the
    W3C validation tool for HTML 4.01 and XHTML 1.0 Transitional.

	[this is a link](https://github.com/kataras/iris) `)

//不应在包含动态数据的处理程序上使用缓存。 在静态内容（即“关于页面”或整个博客网站）上，缓存是一项很好的功能，必须具备

// Cache should not be used on handlers that contain dynamic data.
// Cache is a good and a must-feature on static content, i.e "about page" or for a whole blog site.
func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Get("/", cache.Handler(10*time.Second), writeMarkdown)

	//将其内容保存在第一个请求上，并为它提供服务，而不是重新计算内容
	//10秒钟后，它将清除并重置

	// saves its content on the first request and serves it instead of re-calculating the content.
	// After 10 seconds it will be cleared and reset.

	app.Run(iris.Addr(":8080"))
}

func writeMarkdown(ctx iris.Context) {
	//多次点击浏览器的刷新按钮
	//您每10秒只会看到一次该println

	// tap multiple times the browser's refresh button and you will
	// see this println only once every 10 seconds.
	println("Handler executed. Content refreshed.")

	ctx.Markdown(markdownContents)
}

/*
请注意，`HandleDir`确实会默认使用浏览器的磁盘缓存，因此，在任何HandleDir调用之后注册缓存处理程序，
以实现一种更快的解决方案，即服务器无需跟踪响应即可导航至
 https://github.com/kataras/iris/blob/master/_examples/cache/client-side/main.go

Note that `HandleDir` does use the browser's disk caching by-default
therefore, register the cache handler AFTER any HandleDir calls,
for a faster solution that server doesn't need to keep track of the response
navigate to https://github.com/kataras/iris/blob/master/_examples/cache/client-side/main.go 
*/
```
## 提示
1. 第一次访问，服务器会返回所有信息，当在缓存时间之内请求服务器，将得到最缓存的信息。过期以后将从新在服务生成。
2. 适合于静态页面做缓存