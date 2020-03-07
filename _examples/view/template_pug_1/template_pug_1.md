# go iris 视图pug模板引擎第1个示例
## 目录结构
> 主目录`template_pug_1`
```html
    —— templates
        —— index.pug
    —— main.go
```
## 代码示例
> `templates/index.pug`
```html
doctype html
html(lang="en")
	head
		meta(charset="utf-8")
		title Title
	body
		p ads
		ul
			li The name is {{.Name}}.
			li The age is {{.Age}}.

		each _,_ in .Emails
			div An email is {{.}}

		| {{ with .Jobs }}
			each _,_ in .
				div
				 An employer is {{.Employer}}
				 and the role is {{.Role}}
		| {{ end }}
```
> `main.go`
```golang
//包main显示了一个基于https://github.com/Joker/jade/tree/master/example/actions的哈巴狗动作的示例
//
// Package main shows an example of pug actions based on https://github.com/Joker/jade/tree/master/example/actions
package main

import "github.com/kataras/iris/v12"

type Person struct {
	Name   string
	Age    int
	Emails []string
	Jobs   []*Job
}

type Job struct {
	Employer string
	Role     string
}

func main() {
	app := iris.New()

	tmpl := iris.Pug("./templates", ".pug")
	app.RegisterView(tmpl)

	app.Get("/", index)

	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}

func index(ctx iris.Context) {
	job1 := Job{Employer: "Monash B", Role: "Honorary"}
	job2 := Job{Employer: "Box Hill", Role: "Head of HE"}

	person := Person{
		Name:   "jan",
		Age:    50,
		Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
		Jobs:   []*Job{&job1, &job2},
	}

	ctx.View("index.pug", person)
}
```