# 教程使用Iris和Vue.js的Todo MVC应用程序
## 目录结构
> 主目录`vue-todo-mvc`
```html
    —— src
        —— todo
            —— item.go
            —— service.go
        —— web
            —— controllers
                —— todo_controller.go
            —— public
                —— css
                    —— index
                —— js
                    —— lib
                        —— vue
                    —— app.js
                —— index.html
        —— main.go
```

## Hackernoon文章：https://medium.com/hackernoon/a-todo-mvc-application-using-iris-and-vue-js-5019ff870064

Vue.js是用于使用javascript构建Web应用程序的前端框架。它具有快速的虚拟DOM渲染器

Iris是使用Go编程语言构建Web应用程序的后端框架（免责声明：此处为作者）。 它是目前最快的功能强大的Web框架之一。 我们想用它来服务我们的"todo service"

## 工具
编程语言只是我们的工具，但是我们需要一种安全，快速且“跨平台”的编程语言来为我们的服务提供动力

[Go](https://golang.org) 他是 [快速增长的](https://www.tiobe.com/tiobe-index/) 专为构建简单，快速和可靠的软件而设计的开源编程语言。 看一看 [这里](https://github.com/golang/go/wiki/GoUsers) 哪些伟大的公司使用Go来为其提供服务

### 安装Go编程语言

有关下载和安装Go的广泛信息，请参见 [此处](https://golang.org/dl/).

[![](https://i3.ytimg.com/vi/9x-pG3lvLi0/hqdefault.jpg)](https://youtu.be/9x-pG3lvLi0)

>  [Windows](https://www.youtube.com/watch?v=WT5mTznJBS0) 或 [Mac OS X](https://www.youtube.com/watch?v=5qI8z_lB5Lw) 

> 本文不包含对语言本身的介绍，如果您是新手，建议您为本文添加书签，[学习](https://github.com/golang/go/wiki/Learn)语言的基础知识和 以后再回来

## 依赖关系

过去曾有许多文章导致开发人员不使用web框架，因为它们是无用的和“不好的”,我必须告诉您，没有这样的事情，它总是取决于您要使用的（web）框架。 
在生产环境中，我们没有时间或经验来编写我们想要在应用程序中使用的所有内容的代码，并且如果可以的话，我们确定我们可以做得比别人更好，更安全吗？ 
短期而言：**好的框架对任何开发人员，公司或初创公司都是有用的工具，而“坏的”框架则浪费时间，非常清楚。**

您需要两个依赖项：
1. Vue。用于我们的客户端需求。下载它从[这里](https://vuejs.org/)，最新的v2
2. Iris Web框架，用于我们的服务器端需求。可以找到[这里](https://github.com/kataras/iris)，最新的v12

> 如果你已经安装了，那么只需执行`Go get github.com/kataras/iris/v12@latest`来安装Iris Web框架。
## 开始

如果我们都在同一个页面，那么现在就该学习如何创建一个易于部署和扩展的live todo应用程序了!

我们将使用vue。js todo应用它使用浏览器你可以在vue的[手册](https://vuejs.org/v2/examples/todomvc.html)中找到它的原始版本。

假设您知道%GOPATH%的工作方式，请在%GOPATH%/src目录中创建一个空文件夹，即“ vuejs-todo-mvc”，您将在其中创建这些文件：

- web/public/js/app.js
- web/public/index.html
- todo/item.go
- todo/service.go
- web/controllers/todo_controller.go
- web/main.go

_阅读源代码中的注释，它们可能非常有帮助_

## 代码示例
> `src/todo/item.go`
```go
package todo

type Item struct {
	SessionID string `json:"-"`
	ID        int64  `json:"id,omitempty"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
```
> `src/todo/service.go`
```go
package todo

import (
	"sync"
)

type Service interface {
	Get(owner string) []Item
	Save(owner string, newItems []Item) error
}

type MemoryService struct {
	//键=会话ID，值此会话ID拥有的todo事项列表

	// key = session id, value the list of todo items that this session id has.
	items map[string][]Item
	//由locker保护以进行并发访问

	// protected by locker for concurrent access.
	mu sync.RWMutex
}

func NewMemoryService() *MemoryService {
	return &MemoryService{
		items: make(map[string][]Item, 0),
	}
}

func (s *MemoryService) Get(sessionOwner string) []Item {
	s.mu.RLock()
	items := s.items[sessionOwner]
	s.mu.RUnlock()

	return items
}

func (s *MemoryService) Save(sessionOwner string, newItems []Item) error {
	var prevID int64
	for i := range newItems {
		if newItems[i].ID == 0 {
			newItems[i].ID = prevID
			prevID++
		}
	}

	s.mu.Lock()
	s.items[sessionOwner] = newItems
	s.mu.Unlock()
	return nil
}
```
> `src/web/controllers/todo_controller.go`
```go
package controllers

import (
	"github.com/kataras/iris/v12/_examples/tutorial/vuejs-todo-mvc/src/todo"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/websocket"
)

// TodoController是我们的TODO应用程序的Web控制器

// TodoController is our TODO app's web controller.
type TodoController struct {
	Service todo.Service

	Session *sessions.Session

	NS *websocket.NSConn
}

//在服务器运行之前以及路由和依赖项绑定之前，调用一次BeforeActivation
// 您可以将自定义内容绑定到控制器，添加新方法，添加中间件，向结构或方法添加依赖项等等

// BeforeActivation called once before the server ran, and before
// the routes and dependencies binded.
// You can bind custom things to the controller, add new methods, add middleware,
// add dependencies to the struct or the method(s) and more.
func (c *TodoController) BeforeActivation(b mvc.BeforeActivation) {
	//可以绑定到控制器的函数输入参数
	// （如果有）或struct字段（如果有）：
	b.Dependencies().Add(func(ctx iris.Context) (items []todo.Item) {
		ctx.ReadJSON(&items)
		return
	})
}
// Get处理GET: /todos路线‘

// Get handles the GET: /todos route.
func (c *TodoController) Get() []todo.Item {
	return c.Service.Get(c.Session.ID())
}
// PostItemResponse所有todo items保存动作后将以json形式返回的响应数据

// PostItemResponse the response data that will be returned as json
// after a post save action of all todo items.
type PostItemResponse struct {
	Success bool `json:"success"`
}

var emptyResponse = PostItemResponse{Success: false}
// Post处理POST: /todos路由

// Post handles the POST: /todos route.
func (c *TodoController) Post(newItems []todo.Item) PostItemResponse {
	if err := c.Service.Save(c.Session.ID(), newItems); err != nil {
		return emptyResponse
	}

	return PostItemResponse{Success: true}
}

func (c *TodoController) Save(msg websocket.Message) error {
	id := c.Session.ID()
	c.NS.Conn.Server().Broadcast(nil, websocket.Message{
		Namespace: msg.Namespace,
		Event:     "saved",
		To:        id,
		Body:      websocket.Marshal(c.Service.Get(id)),
	})

	return nil
}
```
> `src/web/public/css/index`
```text
index.css并不是为了减少示例的磁盘空间
而是使用https://unpkg.com/todomvc-app-css@2.0.4/index.css

index.css is not here to reduce the disk space for the examples.
https://unpkg.com/todomvc-app-css@2.0.4/index.css is used instead.
```
> `src/web/public/js/lib/vue`
```text
vue.js并非在此处减少示例的磁盘空间
而是使用https://vuejs.org/js/vue.js代替

vue.js is not here to reduce the disk space for the examples.
Instead https://vuejs.org/js/vue.js is used instead.
```
> `src/web/public/js/app.js`
```js
//完全符合规范的TodoMVC，具有约200条有效的JavaScript行中的iris和基于哈希的路由

// Full spec-compliant TodoMVC with Iris
// and hash-based routing in ~200 effective lines of JavaScript.

var ws;

((async () => {
  const events = {
    todos: {
      saved: function (ns, msg) {
        app.todos = msg.unmarshal()
        //或进行新的http提取

        // or make a new http fetch
        // fetchTodos(function (items) {
        //   app.todos = msg.unmarshal()
        // });
      }
    }
  };

  const conn = await neffos.dial("ws://localhost:8080/todos/sync", events);
  ws = await conn.connect("todos");
})()).catch(console.error);

function fetchTodos(onComplete) {
  axios.get("/todos").then(response => {
    if (response.data === null) {
      return;
    }

    onComplete(response.data);
  });
}

var todoStorage = {
  fetch: function () {
    var todos = [];
    fetchTodos(function (items) {
      for (var i = 0; i < items.length; i++) {
        todos.push(items[i]);
      }
    });
    return todos;
  },
  save: function (todos) {
    axios.post("/todos", JSON.stringify(todos)).then(response => {
      if (!response.data.success) {
        window.alert("saving had a failure");
        return;
      }
      // console.log("send: save");
      ws.emit("save")
    });
  }
}
//可见性过滤器

// visibility filters
var filters = {
  all: function (todos) {
    return todos
  },
  active: function (todos) {
    return todos.filter(function (todo) {
      return !todo.completed
    })
  },
  completed: function (todos) {
    return todos.filter(function (todo) {
      return todo.completed
    })
  }
}
// app Vue实例

// app Vue instance
var app = new Vue({
  // app initial state
  data: {
    todos: todoStorage.fetch(),
    newTodo: '',
    editedTodo: null,
    visibility: 'all'
  },
  // 我们不会使用"watch"，因为它可以与“ hasChanges”和回调之类的字段一起使用，以实现它的真实性，
  // 但是让我们非常简单，因为这只是一个很小的入门
  // we will not use the "watch" as it works with the fields like "hasChanges"
  // and callbacks to make it true but let's keep things very simple as it's just a small getting started.

  // //观看todos更改以保持持久性 |  watch todos change for persistence
  // watch: {
  //   todos: {
  //     handler: function (todos) {
  //       if (app.hasChanges) {
  //         todoStorage.save(todos);
  //         app.hasChanges = false;
  //       }

  //     },
  //     deep: true
  //   }
  // },

  // computed属性 | computed properties
  // http://vuejs.org/guide/computed.html
  computed: {
    filteredTodos: function () {
      return filters[this.visibility](this.todos)
    },
    remaining: function () {
      return filters.active(this.todos).length
    },
    allDone: {
      get: function () {
        return this.remaining === 0
      },
      set: function (value) {
        this.todos.forEach(function (todo) {
          todo.completed = value
        })
        this.notifyChange();
      }
    }
  },

  filters: {
    pluralize: function (n) {
      return n === 1 ? 'item' : 'items'
    }
  },
  //实现数据逻辑的方法
  //注意，这里根本没有DOM操作

  // methods that implement data logic.
  // note there's no DOM manipulation here at all.
  methods: {
    notifyChange: function () {
      todoStorage.save(this.todos)
    },
    addTodo: function () {
      var value = this.newTodo && this.newTodo.trim()
      if (!value) {
        return
      }
      this.todos.push({
        id: this.todos.length + 1, // //仅用于客户端 | just for the client-side.
        title: value,
        completed: false
      })
      this.newTodo = ''
      this.notifyChange();
    },

    completeTodo: function (todo) {
      if (todo.completed) {
        todo.completed = false;
      } else {
        todo.completed = true;
      }
      this.notifyChange();
    },
    removeTodo: function (todo) {
      this.todos.splice(this.todos.indexOf(todo), 1)
      this.notifyChange();
    },

    editTodo: function (todo) {
      this.beforeEditCache = todo.title
      this.editedTodo = todo
    },

    doneEdit: function (todo) {
      if (!this.editedTodo) {
        return
      }
      this.editedTodo = null
      todo.title = todo.title.trim();
      if (!todo.title) {
        this.removeTodo(todo);
      }
      this.notifyChange();
    },

    cancelEdit: function (todo) {
      this.editedTodo = null
      todo.title = this.beforeEditCache
    },

    removeCompleted: function () {
      this.todos = filters.active(this.todos);
      this.notifyChange();
    }
  },
  //一个自定义指令，等待DOM更新后再关注输入字段

  // a custom directive to wait for the DOM to be updated
  // before focusing on the input field.
  // http://vuejs.org/guide/custom-directive.html
  directives: {
    'todo-focus': function (el, binding) {
      if (binding.value) {
        el.focus()
      }
    }
  }
})
//处理路由

// handle routing
function onHashChange() {
  var visibility = window.location.hash.replace(/#\/?/, '')
  if (filters[visibility]) {
    app.visibility = visibility
  } else {
    window.location.hash = ''
    app.visibility = 'all'
  }
}

window.addEventListener('hashchange', onHashChange)
onHashChange()

// mount
app.$mount('.todoapp');
```
> `src/web/public/index.html`
```html
<!doctype html>
<html data-framework="vue">

<head>
  <meta charset="utf-8">
  <title>Iris + Vue.js • TodoMVC</title>
  <link rel="stylesheet" href="https://unpkg.com/todomvc-app-css@2.0.4/index.css">
  <!-- this needs to be loaded before guide's inline scripts -->
  <script src="https://vuejs.org/js/vue.js"></script>
  <!-- $http -->
  <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
  <!-- -->
  <script src="https://unpkg.com/director@1.2.8/build/director.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/neffos.js@latest/dist/neffos.min.js"></script>

  <style>
    [v-cloak] {
      display: none;
    }
  </style>
</head>

<body>
  <section class="todoapp">
    <header class="header">
      <h1>todos</h1>
      <input class="new-todo" autofocus autocomplete="off" placeholder="What needs to be done?" v-model="newTodo"
        @keyup.enter="addTodo">
    </header>
    <section class="main" v-show="todos.length" v-cloak>
      <input class="toggle-all" type="checkbox" v-model="allDone">
      <ul class="todo-list">
        <li v-for="todo in filteredTodos" class="todo" :key="todo.id"
          :class="{ completed: todo.completed, editing: todo == editedTodo }">
          <div class="view">
            <!-- v-model="todo.completed" -->
            <input class="toggle" type="checkbox" @click="completeTodo(todo)">
            <label @dblclick="editTodo(todo)">{{ todo.title }}</label>
            <button class="destroy" @click="removeTodo(todo)"></button>
          </div>
          <input class="edit" type="text" v-model="todo.title" v-todo-focus="todo == editedTodo" @blur="doneEdit(todo)"
            @keyup.enter="doneEdit(todo)" @keyup.esc="cancelEdit(todo)">
        </li>
      </ul>
    </section>
    <footer class="footer" v-show="todos.length" v-cloak>
      <span class="todo-count">
        <strong>{{ remaining }}</strong> {{ remaining | pluralize }} left
      </span>
      <ul class="filters">
        <li>
          <a href="#/all" :class="{ selected: visibility == 'all' }">All</a>
        </li>
        <li>
          <a href="#/active" :class="{ selected: visibility == 'active' }">Active</a>
        </li>
        <li>
          <a href="#/completed" :class="{ selected: visibility == 'completed' }">Completed</a>
        </li>
      </ul>
      <button class="clear-completed" @click="removeCompleted" v-show="todos.length > remaining">
        Clear completed
      </button>
    </footer>
  </section>
  <footer class="info">
    <p>Double-click to edit a todo</p>
  </footer>

  <script src="/js/app.js"></script>
</body>

</html>
```
> `src/web/main.go`
```go
package main

import (
	"strings"

	"github.com/kataras/iris/v12/_examples/tutorial/vuejs-todo-mvc/src/todo"
	"github.com/kataras/iris/v12/_examples/tutorial/vuejs-todo-mvc/src/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/websocket"
)

func main() {
	app := iris.New()
	// 在public目录中提供我们的应用程序，公共文件夹包含客户端vue.js应用程序，
	// 这里不需要任何服务器端模板，
	// 实际上，如果您只使用vue而没有任何后端服务，
	// 则可以 在此行之后停止并启动服务器。

	// serve our app in public, public folder
	// contains the client-side vue.js application,
	// no need for any server-side template here,
	// actually if you're going to just use vue without any
	// back-end services, you can just stop afer this line and start the server.
	app.HandleDir("/", "./public")
	//配置http会话

	// configure the http sessions.
	sess := sessions.New(sessions.Config{
		Cookie: "iris_session",
	})
	//创建一个子路由器并注册http controllers

	// create a sub router and register the http controllers.
	todosRouter := app.Party("/todos")
	//创建针对/todos相对子路径的mvc应用程序

	// create our mvc application targeted to /todos relative sub path.
	todosApp := mvc.New(todosRouter)
	//这里的所有依赖项绑定...

	// any dependencies bindings here...
	todosApp.Register(
		todo.NewMemoryService(),
		sess.Start,
	)

	todosController := new(controllers.TodoController)
	//控制器注册在这里...

	// controllers registration here...
	todosApp.Handle(todosController)
	//为websocket控制器创建一个子mvc应用
	//继承父级的依赖项

	// Create a sub mvc app for websocket controller.
	// Inherit the parent's dependencies.
	todosWebsocketApp := todosApp.Party("/sync")
	todosWebsocketApp.HandleWebsocket(todosController).
		SetNamespace("todos").
		SetEventMatcher(func(methodName string) (string, bool) {
			return strings.ToLower(methodName), true
		})

	websocketServer := websocket.New(websocket.DefaultGorillaUpgrader, todosWebsocketApp)
	idGenerator := func(ctx iris.Context) string {
		id := sess.Start(ctx).ID()
		return id
	}
	todosWebsocketApp.Router.Get("/", websocket.Handler(websocketServer, idGenerator))
	//在http://localhost:8080启动Web服务器

	// start the web server at http://localhost:8080
	app.Run(iris.Addr(":8080"))
}
```
通过从当前路径(%GOPATH%/src/%your_folder%/web/)执行`go run main.go`，运行刚刚创建的Iris Web服务器

```sh
$ go run main.go
Now listening on: http://localhost:8080
Application Started. Press CTRL+C to shut down.
_
```

在http://localhost:8080上打开一个或多个浏览器选项卡，玩得开心！

![](screen.png)

### 下载源代码

整个项目以及您在本文中看到的所有文件都位于：https://github.com/kataras/iris/tree/master/_examples/tutorial/vuejs-todo-mvc

## 参考文献

https://vuejs.org/v2/examples/todomvc.html（使用浏览器的本地存储）

https://github.com/kataras/iris/tree/master/_examples/mvc（mvc示例和功能概述存储库）

## 再一次谢谢你

新年快乐，再次感谢您的耐心配合。）请随时提出任何问题并提供反馈（我是非常活跃的开发人员，因此会在这里听到您的声音！）

别忘了查看我的个人资料和Twitter，我也在那里发布了一些（有用的）信息：）

- https://medium.com/@kataras 
- https://twitter.com/MakisMaropoulos