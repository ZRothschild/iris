# IRIS `seesion`会话加密cookie
## 目录结构
> 主目录`securecookie`
```html
    —— main.go
    —— main_test.go
```
## 代码示例
> `main.go`
```golang
package main

//开发人员可以使用任何库来添加自定义cookie编码器/解码器(encoder/decoder.)
//在此示例中，我们使用gorilla的securecookie软件包：
// $ go get github.com/gorilla/securecookie
// $ go run main.go

// developers can use any library to add a custom cookie encoder/decoder.
// At this example we use the gorilla's securecookie package:
// $ go get github.com/gorilla/securecookie
// $ go run main.go

import (
	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/sessions"

	"github.com/gorilla/securecookie"
)

func newApp() *iris.Application {
	app := iris.New()

	cookieName := "mycustomsessionid"
	// AES仅支持16、24或32字节的密钥大小。
	//您需要提供确切的数据，要么写死或从其他地方提供

	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	hashKey := []byte("the-big-and-secret-fash-key-here")
	blockKey := []byte("lot-secret-of-characters-big-too")
	secureCookie := securecookie.New(hashKey, blockKey)

	mySessions := sessions.New(sessions.Config{
		Cookie:       cookieName,
		Encode:       secureCookie.Encode,
		Decode:       secureCookie.Decode,
		AllowReclaim: true,
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("You should navigate to the /set, /get, /delete, /clear,/destroy instead")
	})
	app.Get("/set", func(ctx iris.Context) {
		//设置会话值
		// set session values
		s := mySessions.Start(ctx)
		s.Set("name", "iris")
		//测试是否在这里设置
		// test if set here
		ctx.Writef("All ok session set to: %s", s.GetString("name"))
	})

	app.Get("/get", func(ctx iris.Context) {
		//以字符串的形式获取特定键，如果找不到，则仅返回一个空字符串
		// get a specific key, as string, if no found returns just an empty string
		s := mySessions.Start(ctx)
		name := s.GetString("name")

		ctx.Writef("The name on the /set was: %s", name)
	})

	app.Get("/delete", func(ctx iris.Context) {
		//删除特定的key
		// delete a specific key
		s := mySessions.Start(ctx)
		s.Delete("name")
	})

	app.Get("/clear", func(ctx iris.Context) {
		//删除所有key
		// removes all entries
		mySessions.Start(ctx).Clear()
	})

	app.Get("/update", func(ctx iris.Context) {
		//更新过期时间用一个新时间
		// updates expire date with a new date
		mySessions.ShiftExpiration(ctx)
	})

	app.Get("/destroy", func(ctx iris.Context) {
		//销毁，删除整个会话数据和cookie
		// destroy, removes the entire session data and cookie
		mySessions.Destroy(ctx)
	})
	//注意销毁：
	//
	//您也可以使用以下命令销毁处理程序外部的会话：
	// mySessions.DestroyByID
	// mySessions.DestroyAll
	
	// Note about destroy:
	//
	// You can destroy a session outside of a handler too, using the:
	// mySessions.DestroyByID
	// mySessions.DestroyAll

	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}
```
> `main_test.go`
```golang
package main

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestSessionsEncodeDecode(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app, httptest.URL("http://example.com"))

	es := e.GET("/set").Expect()
	es.Status(iris.StatusOK)
	es.Cookies().NotEmpty()
	es.Body().Equal("All ok session set to: iris")

	e.GET("/get").Expect().Status(iris.StatusOK).Body().Equal("The name on the /set was: iris")
	// delete and re-get
	e.GET("/delete").Expect().Status(iris.StatusOK)
	e.GET("/get").Expect().Status(iris.StatusOK).Body().Equal("The name on the /set was: ")
	// set, clear and re-get
	e.GET("/set").Expect().Body().Equal("All ok session set to: iris")
	e.GET("/clear").Expect().Status(iris.StatusOK)
	e.GET("/get").Expect().Status(iris.StatusOK).Body().Equal("The name on the /set was: ")
}
```