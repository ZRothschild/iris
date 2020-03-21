# go iris 视图 pug 模板第三示例
## 目录结构
> 主目录`template_pug_3`
```html
    —— templates
        —— index.pug
        —— layout.pug
    —— main.go
```
## 代码示例
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
> `templates/layout.pug`
```html
doctype html
html
  head
    block title
      title Default title
  body
    block content
```
> `main.go`
```golang
package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()

	tmpl := iris.Pug("./templates", ".pug")

	app.RegisterView(tmpl)

	app.Get("/", index)

	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}

func index(ctx iris.Context) {
	ctx.View("index.pug")
}
```