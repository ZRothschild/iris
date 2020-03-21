# `go iris caddy`使用示例二
## 目录结构
> 主目录`server2`
```html
    —— Caddyfile
    —— main.go
```
## 代码示例
> `Caddyfile`
```editorconfig
example.com {
	header / Server "Iris"
	proxy / example.com:9091 # localhost:9091
}

api.example.com {
	header / Server "Iris"
	proxy / api.example.com:9092 # localhost:9092
}
```
> `main.go`
```golang
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type postValue func(string) string

func main() {
	app := iris.New()

	mvc.New(app.Party("/user")).Register(
		func(ctx iris.Context) postValue {
			return ctx.PostValue
		}).Handle(new(UserController))

	// GET http://localhost:9092/user
	// GET http://localhost:9092/user/42
	// POST http://localhost:9092/user
	// PUT http://localhost:9092/user/42
	// DELETE http://localhost:9092/user/42
	// GET http://localhost:9092/user/followers/42
	app.Run(iris.Addr(":9092"))
}

// UserController是我们的用户示例控制器。

// UserController is our user example controller.
type UserController struct{}

//Get处理 GET类型 /user请求

// Get handles GET /user
func (c *UserController) Get() string {
	return "Select all users"
}

//User是我们的测试用户模型，这里没有什么大不了的。

// User is our test User model, nothing tremendous here.
type User struct{ ID int64 }

// GetBy处理GET /user/42等于.Get("/user/{id:long}")

// GetBy handles GET /user/42, equal to .Get("/user/{id:long}")
func (c *UserController) GetBy(id int64) User {
	//通过ID == $id选择User

	// Select User by ID == $id.
	return User{id}
}
// Post处理POST /user

// Post handles POST /user
func (c *UserController) Post(post postValue) string {
	username := post("username")
	return "Create by user with username: " + username
}

// PutBy处理PUT /user/42

// PutBy handles PUT /user/42
func (c *UserController) PutBy(id int) string {
	//通过ID == $id更新User
	// Update user by ID == $id
	return "User updated"
}

// DeleteBy处理DELETE /user/42

// DeleteBy handles DELETE /user/42
func (c *UserController) DeleteBy(id int) bool {
	//通过ID ==$id删除User
	//
	//当boolean时，则为true = iris.StatusOK，false = iris.StatusNotFound

	// Delete user by ID == %id
	//
	// when boolean then true = iris.StatusOK, false = iris.StatusNotFound
	return true
}

// GetFollowersBy处理GET /user/followers/42

// GetFollowersBy handles GET /user/followers/42
func (c *UserController) GetFollowersBy(id int) []User {
	//通过用户IID == $id选择所有关注者

	// Select all followers by user ID == $id
	return []User{ /* ... */ }
}
```