# go iris 视图写入`io.Write`
## 目录结构
> 主目录`write-to`
```html
    —— views
        —— email
            —— simple.html
        —— shared
            —— email.html
    —— main.go
```
## 代码示例
> `views/email/simple.html`
```html
{{.Body}}
```
> `views/shared/email.html`
```html
<h1>{{.Title}}</h1>
<p class="body">
    {{yield}}
</p>

<a href="{{.RefLink}}" target="_new">{{.RefTitle}}</a>
```
> `main.go`
```golang
package main

import (
	"os"

	"github.com/kataras/iris/v12"
)

type mailData struct {
	Title    string
	Body     string
	RefTitle string
	RefLink  string
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.RegisterView(iris.HTML("./views", ".html"))

	//您需要在使用`app.View`函数之前手动调用`app.Build`，因此模板是在该状态下构建的

	// you need to call `app.Build` manually before using the `app.View` func,
	// so templates are built in that state.
	app.Build()
	//或使用字符串缓冲的编写器来利用其正文发送电子邮件以发送电子邮件，
	//您可以使用https://github.com/kataras/go-mailer或您喜欢的任何其他第三方程序包

	// Or a string-buffered writer to use its body to send an e-mail
	// for sending e-mails you can use the https://github.com/kataras/go-mailer
	// or any other third-party package you like.
	
	//模板的解析结果将被写入该编写器
	//
	// The template's parsed result will be written to that writer.
	writer := os.Stdout
	err := app.View(writer, "email/simple.html", "shared/email.html", mailData{
		Title:    "This is my e-mail title",
		Body:     "This is my e-mail body",
		RefTitle: "Iris web framework",
		RefLink:  "https://iris-go.com",
	})
	if err != nil {
		app.Logger().Errorf("error from app.View: %v", err)
	}

	app.Run(iris.Addr(":8080"))
}
```