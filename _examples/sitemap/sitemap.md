# `sitemap`网站地图使用介绍
## 网站地图介绍
网站地图是一个网站的缩影，包含网站的内容地址，是根据网站的结构、框架、内容，生成的导航文件。网站地图分为三种文件格式：xml格式、html格式以及txt格式。xml格式和txt格式一般用于搜索引擎，为搜索引擎蜘蛛程序提供便利的入口到你的网站所有网页；html格式网站地图可以作为一个网页展示给访客，方便用户查看网站内容。本文主要讲解xml格式的网站地图相关内容和生成方法。

Sitemap 0.90 是依据创意公用授权-相同方式共享的条款提供的，并被广泛采用，受 Google、Yahoo! 和 Microsoft 在内的众多厂商的支持
## 目录结构
> 主目录`sitemap`
```html
    —— main.go
    —— main_test.go
```
## 代码示例
> `main.go`
```go
package main

import (
	"time"

	"github.com/kataras/iris/v12"
)

const startURL = "http://localhost:8080"

func main() {
	app := newApp()

	// http://localhost:8080/sitemap.xml
	//仅列出在线GET静态路由
	// Lists only online GET static routes.
	//
	// Reference/参考 : https://www.sitemaps.org/protocol.html
	app.Run(iris.Addr(":8080"), iris.WithSitemap(startURL))
}

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug")

	lastModified, _ := time.Parse("2006-01-02T15:04:05-07:00", "2019-12-13T21:50:33+02:00")
	app.Get("/home", handler).SetLastMod(lastModified).SetChangeFreq("hourly").SetPriority(1)
	app.Get("/articles", handler).SetChangeFreq("daily")
	app.Get("/path1", handler)
	app.Get("/path2", handler)

	app.Post("/this-should-not-be-listed", handler)
	app.Get("/this/{myparam}/should/not/be/listed", handler)
	app.Get("/this-should-not-be-listed-offline", handler).SetStatusOffline()

	return app
}

func handler(ctx iris.Context) { ctx.WriteString(ctx.Path()) }
```

> `main_test.go`
```go
package main

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestSitemap(t *testing.T) {
	const expectedFullSitemapXML = `<?xml version="1.0" encoding="utf-8" standalone="yes"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>http://localhost:8080/home</loc><lastmod>2019-12-13T21:50:33+02:00</lastmod><changefreq>hourly</changefreq><priority>1</priority></url><url><loc>http://localhost:8080/articles</loc><changefreq>daily</changefreq></url><url><loc>http://localhost:8080/path1</loc></url><url><loc>http://localhost:8080/path2</loc></url></urlset>`

	app := newApp()
	app.Configure(iris.WithSitemap(startURL))

	e := httptest.New(t, app)
	e.GET("/sitemap.xml").Expect().Status(httptest.StatusOK).Body().Equal(expectedFullSitemapXML)
}
```