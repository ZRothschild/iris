# go iris 视图pug模板引擎第2个示例
## 目录结构
> 主目录`template_pug_2`
```html
    —— templates
        —— footer.pug
        —— header.pug
        —— index.pug
    —— main.go
```
## 代码示例
> `templates/footer.pug`
```html
#footer
  p Copyright (c) foobar
```
> `templates/header.pug`
```html
head
  title My Site
  <!-- script(src='/javascripts/jquery.js')
  script(src='/javascripts/app.js') -->
```
> `templates/index.pug`
```html
doctype html
html
  include templates/header.pug
  body
    h1 My Site
    p {{ bold "Welcome to my super lame site."}}
    include templates/footer.pug
```
> `main.go`
```golang
package main

import (
	"html/template"

	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	tmpl := iris.Pug("./templates", ".pug")
	//根据每个请求重新加载模板（开发模式）
	tmpl.Reload(true)                                            // reload templates on each request (development mode)
	//在此处添加模板功能
	tmpl.AddFunc("bold", func(s string) (template.HTML, error) { // add your template func here.
		return template.HTML("<b>" + s + "</b>"), nil
	})

	app.RegisterView(tmpl)

	app.Get("/", index)

	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}

func index(ctx iris.Context) {
	ctx.View("index.pug")
}
```