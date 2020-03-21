# go iris 视图 jet 模板第零示例
## 目录结构
> 主目录`template_jet_0`
```html
    —— views
        —— layouts
            —— application.jet
        —— todos
            —— index.jet
            —— show.jet
    —— main.go
```
## 代码示例
> `views/layouts/application.jet`
```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>{{ isset(title) ? title : "" }}</title>
  </head>
  <body>
    {{block documentBody()}}{{end}}
  </body>
</html>
```
> `views/todos/index.jet`
```html
{{extends "layouts/application.jet"}}

{{block button(label, href="javascript:void(0)")}}
  <a href="{{ href }}">{{ label }}</a>
{{end}}

{{block ul()}}
  <ul>
    {{yield content}}
  </ul>
{{end}}

{{block documentBody()}}
  <h1>List of TODOs</h1>

  {{if isset(showingAllDone) && showingAllDone}}
    <p>Showing only TODOs that are done</p>
  {{else}}
    <p><a href="/all-done">Show only TODOs that are done</a></p>
  {{end}}

  {{yield ul() content}}
    {{range id, value := .}}
      <li {{if value.Done}}style="color:red;text-decoration: line-through;"{{end}}>
        <a href="/todo?id={{ id }}">{{ value.Text }}</a>
        {{yield button(label="UP", href="/update/?id="+base64(id))}} - {{yield button(href="/delete/?id="+id, label="DL")}}
      </li>
    {{end}}
  {{end}}
{{end}}
```
> `views/todos/show.jet`
```html
{{extends "layouts/application.jet"}}

{{block documentBody()}}
  <h1>Show TODO</h1>
  <p>This uses a custom renderer by implementing the jet.Renderer (or view.JetRenderer) interface.
  <p>
    {{.}}
  </p>
{{end}}
```
> `main.go`
```golang
//包main展示了如何使用Iris内置的Jet视图引擎轻松使用jet模板解析器
//此示例是https://github.com/CloudyKit/jet/tree/master/examples/todos的自定义分支，因此您可以并排注意到它们之间的差异
//
// Package main shows how to use jet template parser with ease using the Iris built-in Jet view engine.
// This example is customized fork of https://github.com/CloudyKit/jet/tree/master/examples/todos, so you can
// notice the differences side by side.
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/view"
)

type tTODO struct {
	Text string
	Done bool
}

type doneTODOs struct {
	list map[string]*tTODO
	keys []string
	len  int
	i    int
}

func (dt *doneTODOs) New(todos map[string]*tTODO) *doneTODOs {
	dt.len = len(todos)
	for k := range todos {
		dt.keys = append(dt.keys, k)
	}
	dt.list = todos
	return dt
}
// Range满足jet.Ranger接口，即使列表中包含未完成的TODO，也仅返回已完成的TODO

// Range satisfies the jet.Ranger interface and only returns TODOs that are done,
// even when the list contains TODOs that are not done.
func (dt *doneTODOs) Range() (reflect.Value, reflect.Value, bool) {
	for dt.i < dt.len {
		key := dt.keys[dt.i]
		dt.i++
		if dt.list[key].Done {
			return reflect.ValueOf(key), reflect.ValueOf(dt.list[key]), false
		}
	}
	return reflect.Value{}, reflect.Value{}, true
}

func (dt *doneTODOs) Render(r *view.JetRuntime) {
	r.Write([]byte(fmt.Sprintf("custom renderer")))
}
//渲染实现jet.Renderer接口

// Render implements jet.Renderer interface
func (t *tTODO) Render(r *view.JetRuntime) {
	done := "yes"
	if !t.Done {
		done = "no"
	}
	r.Write([]byte(fmt.Sprintf("TODO: %s (done: %s)", t.Text, done)))
}

func main() {
	//输入别名：
	//
	// Type aliases:
	//
	// view.JetRuntimeVars = jet.VarMap
	// view.JetRuntime = jet.Runtime
	// view.JetArguments = jet.Arguments
	//
	// Iris还使您能够通过以下方式放置来自中间件的运行时变量：
	//
	// Iris also gives you the ability to put runtime variables
	// from middlewares as well, by:
	//
	// view.AddJetRuntimeVars(ctx, vars)
	// or tmpl.AddRuntimeVars(ctx, vars)
	//
	//当自定义jet.Ranger和自定义jet.Renderer根本无法工作时，Iris Jet修复了该问题，希望这是Jet解析器本身的临时问题，
	//作者将尽快修复，以便我们删除“hacks” 我们已经让那些工具按预期工作。
	//
	// The Iris Jet fixes the issue when custom jet.Ranger and custom jet.Renderer are not actually work at all,
	// hope that this is a temp issue on the jet parser itself and authors will fix it soon
	// so we can remove the "hacks" we've putted for those to work as expected.

	app := iris.New()
	tmpl := iris.Jet("./views", ".jet") // <--
	tmpl.Reload(true)                   // 在生产中请删除 | remove in production.
	tmpl.AddFunc("base64", func(a view.JetArguments) reflect.Value {
		a.RequireNumOfArguments("base64", 1, 1)

		buffer := bytes.NewBuffer(nil)
		fmt.Fprint(buffer, a.Get(0))

		return reflect.ValueOf(base64.URLEncoding.EncodeToString(buffer.Bytes()))
	})
	app.RegisterView(tmpl) // <--

	todos := map[string]*tTODO{
		"example-todo-1": {Text: "Add an show todo page to the example project", Done: true},
		"example-todo-2": {Text: "Add an add todo page to the example project"},
		"example-todo-3": {Text: "Add an update todo page to the example project"},
		"example-todo-4": {Text: "Add an delete todo page to the example project", Done: true},
	}

	app.Get("/", func(ctx iris.Context) {
		ctx.View("todos/index.jet", todos) // <--
		//注意，如果logger级别允许，则`ctx.View`已经记录了该错误并返回了错误

		// Note that the `ctx.View` already logs the error if logger level is allowing it and returns the error.
	})

	app.Get("/todo", func(ctx iris.Context) {
		id := ctx.URLParam("id")
		todo, ok := todos[id]
		if !ok {
			ctx.Redirect("/")
			return
		}

		ctx.View("todos/show.jet", todo)
	})
	app.Get("/all-done", func(ctx iris.Context) {
		// vars := make(view.JetRuntimeVars)
		// vars.Set("showingAllDone", true)
		// vars.Set("title", "Todos - All Done")
		// view.AddJetRuntimeVars(ctx, vars)
		// ctx.View("todos/index.jet", (&doneTODOs{}).New(todos))
		//
		// OR

		ctx.ViewData("showingAllDone", true)
		ctx.ViewData("title", "Todos - All Done")

		//使用ctx.ViewData（“ _ jet”，jetData）
		//如果要用作中间件
		//预先设置值，甚至以后再从另一个下一个中间件进行更改。
		// ctx.ViewData("_jet", (&doneTODOs{}).New(todos))
		//和ctx.View("todos/index.jet")
		// 要么
		
		// Use ctx.ViewData("_jet", jetData)
		// if using as middleware and you want
		// to pre-set the value or even change it later on from another next middleware.
		// ctx.ViewData("_jet", (&doneTODOs{}).New(todos))
		// and ctx.View("todos/index.jet")
		// OR
		ctx.View("todos/index.jet", (&doneTODOs{}).New(todos))
	})

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	} else if !strings.HasPrefix(":", port) {
		port = ":" + port
	}

	app.Run(iris.Addr(port))
}
```