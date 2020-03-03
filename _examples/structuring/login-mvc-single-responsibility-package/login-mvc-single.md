# iris 登录演示 (MVC 使用独立包组织)，mvc架构简单登录例子
## 目录结构
> 主目录`login-mvc-single-responsibility-package`
```html
    —— public
        —— css
            —— site.css
    —— user
        —— auth.go
        —— controller.go
        —— datasource.go
        —— model.go
    —— views
        —— shared
            —— error.html
            —— layout.html
        —— user
            —— login.html
            —— me.html
            —— notfound.html
            —— register.html
    —— main.go
```
## 项目目录图片
![iris mvc架构简单登录例子](./folder_structure.png)
## 代码示例
> `public/css/site.css`
```css
/* Bordered form */
form {
    border: 3px solid #f1f1f1;
}

/* Full-width inputs */
input[type=text], input[type=password] {
    width: 100%;
    padding: 12px 20px;
    margin: 8px 0;
    display: inline-block;
    border: 1px solid #ccc;
    box-sizing: border-box;
}

/* Set a style for all buttons */
button {
    background-color: #4CAF50;
    color: white;
    padding: 14px 20px;
    margin: 8px 0;
    border: none;
    cursor: pointer;
    width: 100%;
}

/* Add a hover effect for buttons */
button:hover {
    opacity: 0.8;
}

/* Extra style for the cancel button (red) */
.cancelbtn {
    width: auto;
    padding: 10px 18px;
    background-color: #f44336;
}

/* Center the container */

/* Add padding to containers */
.container {
    padding: 16px;
}

/* The "Forgot password" text */
span.psw {
    float: right;
    padding-top: 16px;
}

/* Change styles for span and cancel button on extra small screens */
@media screen and (max-width: 300px) {
    span.psw {
        display: block;
        float: none;
    }
    .cancelbtn {
        width: 100%;
    }
}
```
> `user/auth.go`
```go
package user

import (
	"errors"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

const sessionIDKey = "UserID"

// paths 路径
var (
	PathLogin  = mvc.Response{Path: "/user/login"}
	PathLogout = mvc.Response{Path: "/user/logout"}
)
// AuthController是用户身份验证控制器，是一个自定义共享控制器。

// AuthController is the user authentication controller, a custom shared controller.
type AuthController struct {
	//如果结构体依赖此上下文自动绑定，
	//在此控制器中，我们不会使用mvc样式来做所有事情，
	//但这都不是其功能的30％。
	// Ctx iris.Context

	// context is auto-binded if struct depends on this,
	// in this controller we don't we do everything with mvc-style,
	// and that's neither the 30% of its features.
	// Ctx iris.Context

	Source  *DataSource
	Session *sessions.Session
	//整个控制器处于请求范围内，因为我们已经依赖于Session，因此
	//这对于每个新的传入请求都是新的，BeginRequest根据会话设置它

	// the whole controller is request-scoped because we already depend on Session, so
	// this will be new for each new incoming request, BeginRequest sets that based on the session.
	UserID int64
}
// BeginRequest将登录状态保存到上下文中，即用户ID

// BeginRequest saves login state to the context, the user id.
func (c *AuthController) BeginRequest(ctx iris.Context) {
	c.UserID, _ = c.Session.GetInt64(sessionIDKey)
}
// EndRequest在这里只是为了完成BaseController
//以告知iris在main方法之前调用`BeginRequest`

// EndRequest is here just to complete the BaseController
// in order to be tell iris to call the `BeginRequest` before the main method.
func (c *AuthController) EndRequest(ctx iris.Context) {}

func (c *AuthController) fireError(err error) mvc.View {
	return mvc.View{
		Code: iris.StatusBadRequest,
		Name: "shared/error.html",
		Data: iris.Map{"Title": "User Error", "Message": strings.ToUpper(err.Error())},
	}
}

func (c *AuthController) redirectTo(id int64) mvc.Response {
	return mvc.Response{Path: "/user/" + strconv.Itoa(int(id))}
}

func (c *AuthController) createOrUpdate(firstname, username, password string) (user Model, err error) {
	username = strings.Trim(username, " ")
	if username == "" || password == "" || firstname == "" {
		return user, errors.New("empty firstname, username or/and password")
	}

	userToInsert := Model{
		Firstname: firstname,
		Username:  username,
		password:  password,
	} // password is hashed by the Source. 密码由来源散列

	newUser, err := c.Source.InsertOrUpdate(userToInsert)
	if err != nil {
		return user, err
	}

	return newUser, nil
}

func (c *AuthController) isLoggedIn() bool {
	//我们不按会话搜索，我们有用户ID
	//已经由`BeginRequest`中间件提供

	// we don't search by session, we have the user id
	// already by the `BeginRequest` middleware.
	return c.UserID > 0
}

func (c *AuthController) verify(username, password string) (user Model, err error) {
	if username == "" || password == "" {
		return user, errors.New("please fill both username and password fields")
	}

	u, found := c.Source.GetByUsername(username)
	if !found {
		//如果找不到使用该用户名的用户。

		// if user found with that username not found at all.
		return user, errors.New("user with that username does not exist")
	}

	if ok, err := ValidatePassword(password, u.HashedPassword); err != nil || !ok {
		//如果找到了用户但发生了错误或密码无效。

		// if user found but an error occurred or the password is not valid.
		return user, errors.New("please try to login with valid credentials")
	}

	return u, nil
}
//如果已登录，则销毁会话
//并重定向到登录页面
//否则重定向到注册页面。

// if logged in then destroy the session
// and redirect to the login page
// otherwise redirect to the registration page.
func (c *AuthController) logout() mvc.Response {
	if c.isLoggedIn() {
		c.Session.Destroy()
	}
	return PathLogin
}
```
> `user/controller.go`
```go
package user

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

var (
	//关于代码：iris.StatusSeeOther->
	//当从POST重定向到GET请求时，您-应该-使用此HTTP状态代码，
	//但是，如果您有一些（复杂的）替代方法
	//在线搜索，甚至搜索HTTP RFC
	//另请参见 RFC 7231

	// About Code: iris.StatusSeeOther ->
	// When redirecting from POST to GET request you -should- use this HTTP status code,
	// however there're some (complicated) alternatives if you
	// search online or even the HTTP RFC.
	// "See Other" RFC 7231
	pathMyProfile = mvc.Response{Path: "/user/me", Code: iris.StatusSeeOther}
	pathRegister  = mvc.Response{Path: "/user/register"}
)
//控制器负责处理以下请求：

// Controller is responsible to handle the following requests:
// GET  			/user/register
// POST 			/user/register
// GET 				/user/login
// POST 			/user/login
// GET 				/user/me
// GET				/user/{id:long} | long是一种新的参数类型，它是int64
// GET				/user/{id:long} | long is a new param type, it's the int64.
//所有HTTP方法都试用/user/logout
// All HTTP Methods /user/logout
type Controller struct {
	AuthController
}
//
type formValue func(string) string

//在服务器启动之前和控制器注册之前调用一次BeforeActivation，
// 在这里您可以向该控制器添加依赖项，并且只能添加main调用者可以跳过的依赖项。

// BeforeActivation called once before the server start
// and before the controller's registration, here you can add
// dependencies, to this controller and only, that the main caller may skip.
func (c *Controller) BeforeActivation(b mvc.BeforeActivation) {
	//也绑定上下文的`FormValue`，以便在控制器或其方法的输入参数上也可以接受（也具有NEW功能）。

	// bind the context's `FormValue` as well in order to be
	// acceptable on the controller or its methods' input arguments (NEW feature as well).
	b.Dependencies().Add(func(ctx iris.Context) formValue { return ctx.FormValue })
}

type page struct {
	Title string
}
// GetRegister处理GET:/user/register.
// mvc.Result可以接受任何包含`Dispatch(ctx iris.Context)`方法的结构。
// mvc.Response和mvc.View均为mvc.Result。

// GetRegister handles GET:/user/register.
// mvc.Result can accept any struct which contains a `Dispatch(ctx iris.Context)` method.
// Both mvc.Response and mvc.View are mvc.Result.
func (c *Controller) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		return c.logout()
	}
	//您可以将其用作变量来节约时间，
	//这对您来说是挺不错：)

	// You could just use it as a variable to win some time in serve-time,
	// this is an exersise for you :)
	return mvc.View{
		Name: pathRegister.Path + ".html",
		Data: page{"User Registration"},
	}
}
// PostRegister处理POST:/user/register.

// PostRegister handles POST:/user/register.
func (c *Controller) PostRegister(form formValue) mvc.Result {
	//我们可以使用`c.Ctx.ReadForm`或如下一一读取值

	// we can either use the `c.Ctx.ReadForm` or read values one by one.
	var (
		firstname = form("firstname")
		username  = form("username")
		password  = form("password")
	)

	user, err := c.createOrUpdate(firstname, username, password)
	if err != nil {
		return c.fireError(err)
	}
	//设置会话值从未如此简单

	// setting a session value was never easier.
	c.Session.Set(sessionIDKey, user.ID)
	//成功，这里无事可做，只需重定向到/user/me

	// succeed, nothing more to do here, just redirect to the /user/me.
	return pathMyProfile
}
//通过这些静态视图，
//您可以使用变量-在服务器启动之前进行初始化
//这样您就可以赢得一些服务时间。
//您也可以在其他地方做，但是我让他们为您使用，
//基本上，您可以通过下面的内容来了解

// with these static views,
// you can use variables-- that are initialized before server start
// so you can win some time on serving.
// You can do it else where as well but I let them as pracise for you,
// essentially you can understand by just looking below.
var userLoginView = mvc.View{
	Name: PathLogin.Path + ".html",
	Data: page{"User Login"},
}
// GetLogin处理GET:/user/login.

// GetLogin handles GET:/user/login.
func (c *Controller) GetLogin() mvc.Result {
	if c.isLoggedIn() {
		return c.logout()
	}
	return userLoginView
}
// PostLogin处理POST:/user/login

// PostLogin handles POST:/user/login.
func (c *Controller) PostLogin(form formValue) mvc.Result {
	var (
		username = form("username")
		password = form("password")
	)

	user, err := c.verify(username, password)
	if err != nil {
		return c.fireError(err)
	}

	c.Session.Set(sessionIDKey, user.ID)
	return pathMyProfile
}
// AnyLogout处理路径/user/logout.上的任何http方法

// AnyLogout handles any method on path /user/logout.
func (c *Controller) AnyLogout() {
	c.logout()
}
// GetMe处理GET:/user/me

// GetMe handles GET:/user/me.
func (c *Controller) GetMe() mvc.Result {
	id, err := c.Session.GetInt64(sessionIDKey)
	if err != nil || id <= 0 {
		//如果尚未登录，请重定向至登录

		// when not already logged in, redirect to login.
		return PathLogin
	}

	u, found := c.Source.GetByID(id)
	if !found {
		//如果会话存在，但由于某种原因该用户不存在于“数据库”中
		//然后注销他并重定向到注册页面。

		// if the  session exists but for some reason the user doesn't exist in the "database"
		// then logout him and redirect to the register page.
		return c.logout()
	}
	//设置模型并渲染视图模板

	// set the model and render the view template.
	return mvc.View{
		Name: pathMyProfile.Path + ".html",
		Data: iris.Map{
			"Title": "Profile of " + u.Username,
			"User":  u,
		},
	}
}

func (c *Controller) renderNotFound(id int64) mvc.View {
	return mvc.View{
		Code: iris.StatusNotFound,
		Name: "user/notfound.html",
		Data: iris.Map{
			"Title": "User Not Found",
			"ID":    id,
		},
	}
}
//Dispatch完成`mvc.Result`接口
//是为了能够返回Model的类型
//像mvc.Result可以返回
//如果Dispatch该函数不存在，则 我们应该将输出结果显式设置为该模型或接口{}。

// Dispatch completes the `mvc.Result` interface
// in order to be able to return a type of `Model`
// as mvc.Result.
// If this function didn't exist then
// we should explicit set the output result to that Model or to an interface{}.
func (u Model) Dispatch(ctx iris.Context) {
	ctx.JSON(u)
}
// GetBy处理GET:/user/{id:long}

// GetBy handles GET:/user/{id:long},
// i.e http://localhost:8080/user/1
func (c *Controller) GetBy(userID int64) mvc.Result {
	//我们有/user/{id}
	//获取并呈现用户json

	// we have /user/{id}
	// fetch and render user json.
	user, found := c.Source.GetByID(userID)
	if !found {
		// not user found with that ID.
		return c.renderNotFound(userID)
	}

	//问：模型如何作为mvc.Result返回？
	// A：我之前在一些评论和文档中告诉过您，任何具有`Dispatch(ctx iris.Context)`结构的结构都可以作为mvc.Result返回
	// 因此我们可以 在同一方法中组合多种类型的结果
	// 例如，在这里，我们返回一个mvc.View来呈现未找到的自定义模板
	// 可以是通过Dispatch将模型作为JSON返回的用户

	// Q: how the hell Model can be return as mvc.Result?
	// A: I told you before on some comments and the docs,
	// any struct that has a `Dispatch(ctx iris.Context)`
	// can be returned as an mvc.Result(see ~20 lines above),
	// therefore we are able to combine many type of results in the same method.
	// For example, here, we return either an mvc.View to render a not found custom template
	// either a user which returns the Model as JSON via its Dispatch.

	// //如果`GetBy`的输出结果是该结构的类型或接口{}，
	// 并且iris也将使用JSON进行渲染，则也可以仅返回不是mvc.Result的结构值
	// 但是在这里，我们可以 如果没有完成`Dispatch`函数，就不要这样做，
	// 因为我们可能会返回一个mvc.View，它是一个mvc.Resultl类型

	// We could also return just a struct value that is not an mvc.Result,
	// if the output result of the `GetBy` was that struct's type or an interface{}
	// and iris would render that with JSON as well, but here we can't do that without complete the `Dispatch`
	// function, because we may return an mvc.View which is an mvc.Result.
	return user
}
```
> `user/datasource.go`
```go
package user

import (
	"errors"
	"sync"
	"time"
)

// IDGenerator将是我们的用户ID生成器，
// 但是在这里，我们按用户ID保持用户的顺序，
// 因此我们将使用可以轻松写入浏览器的数字来从REST API中获取结果
// var IDGenerator = func() string {
// 	return uuid.NewV4().String()
// }

// IDGenerator would be our user ID generator
// but here we keep the order of users by their IDs
// so we will use numbers that can be easly written
// to the browser to get results back from the REST API.
// var IDGenerator = func() string {
// 	return uuid.NewV4().String()
// }

// DataSource是我们的数据存储示例

// DataSource is our data store example.
type DataSource struct {
	Users map[int64]Model
	mu    sync.RWMutex
}
// NewDataSource返回一个新的用户数据源

// NewDataSource returns a new user data source.
func NewDataSource() *DataSource {
	return &DataSource{
		Users: make(map[int64]Model),
	}
}

// GetBy接收一个查询函数，
// 该函数针对我们虚构数据库中的每个单个用户模型触发
// 当该函数返回true时，它将停止迭代

// GetBy receives a query function
// which is fired for every single user model inside
// our imaginary database.
// When that function returns true then it stops the iteration.

//返回查询的返回最后一个已知布尔值
// 和最后一个已知用户模型，以帮助调用者减少loc。

// It returns the query's return last known boolean value
// and the last known user model
// to help callers to reduce the loc.

//但要小心，调用者应始终检查“找到”的内容，
//因为它可能为假，但用户模型中实际上具有真实数据。

// But be carefully, the caller should always check for the "found"
// because it may be false but the user model has actually real data inside it.

//实际上，这是我想到的一个简单但非常聪明的原型函数，
//此后一直在各处使用，希望您也发现它非常有用

// It's actually a simple but very clever prototype function
// I'm think of and using everywhere since then,
// hope you find it very useful too.
func (d *DataSource) GetBy(query func(Model) bool) (user Model, found bool) {
	d.mu.RLock()
	for _, user = range d.Users {
		found = query(user)
		if found {
			break
		}
	}
	d.mu.RUnlock()
	return
}
// GetByID根据其ID返回用户模型

// GetByID returns a user model based on its ID.
func (d *DataSource) GetByID(id int64) (Model, bool) {
	return d.GetBy(func(u Model) bool {
		return u.ID == id
	})
}
// GetByUsername返回基于用户名的用户模型

// GetByUsername returns a user model based on the Username.
func (d *DataSource) GetByUsername(username string) (Model, bool) {
	return d.GetBy(func(u Model) bool {
		return u.Username == username
	})
}

func (d *DataSource) getLastID() (lastID int64) {
	d.mu.RLock()
	for id := range d.Users {
		if id > lastID {
			lastID = id
		}
	}
	d.mu.RUnlock()

	return lastID
}
// InsertOrUpdate将用户添加或更新到内存存储

// InsertOrUpdate adds or updates a user to the (memory) storage.
func (d *DataSource) InsertOrUpdate(user Model) (Model, error) {
	//无论我们将更新update和insert动作的密码哈希值

	// no matter what we will update the password hash
	// for both update and insert actions.
	hashedPassword, err := GeneratePassword(user.password)
	if err != nil {
		return user, err
	}
	user.HashedPassword = hashedPassword

	// update
	if id := user.ID; id > 0 {
		_, found := d.GetByID(id)
		if !found {
			return user, errors.New("ID should be zero or a valid one that maps to an existing User")
		}
		d.mu.Lock()
		d.Users[id] = user
		d.mu.Unlock()
		return user, nil
	}

	// insert
	id := d.getLastID() + 1
	user.ID = id
	d.mu.Lock()
	user.CreatedAt = time.Now()
	d.Users[id] = user
	d.mu.Unlock()

	return user, nil
}
```
> `user/model.go`
```go
package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

//模型是我们的用户示例模型

// Model is our User example model.
type Model struct {
	ID        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Username  string `json:"username"`

	//密码是客户端提供的密码，不会存储在服务器中的任何地方
	//它仅用于注册和更新密码之类的操作，
	//因为我们接受`DataSource#InsertOrUpdate`函数中的Model实例

	// password is the client-given password
	// which will not be stored anywhere in the server.
	// It's here only for actions like registration and update password,
	// because we caccept a Model instance
	// inside the `DataSource#InsertOrUpdate` function.
	password       string
	HashedPassword []byte    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

// GeneratePassword将根据用户输入为我们生成一个哈希密码

// GeneratePassword will generate a hashed password for us based on the
// user's input.
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}
// ValidatePassword将检查密码是否匹配

// ValidatePassword will check if passwords are matched.
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}
```
> `views/shared/error.html`
```html
<h1>Error.</h1>
<h2>An error occurred while processing your request.</h2>

<h3>{{.Message}}</h3>
```
> `views/shared/layout.html`
```html
<html>

<head>
    <title>{{.Title}}</title>
    <link rel="stylesheet" type="text/css" href="/public/css/site.css" />
</head>

<body>
    {{ yield }}
</body>

</html>
```
> `views/user/login.html`
```html
<form action="/user/login" method="POST">
    <div class="container">
        <label><b>Username</b></label>
        <input type="text" placeholder="Enter Username" name="username" required>

        <label><b>Password</b></label>
        <input type="password" placeholder="Enter Password" name="password" required>

        <button type="submit">Login</button>
    </div>
</form>
```
> `views/user/me.html`
```html
<p>
    Welcome back <strong>{{.User.Firstname}}</strong>!
</p>
```
> `views/user/notfound.html`
```html
<p>
    User with ID <strong>{{.ID}}</strong> does not exist.
</p>
```
> `views/user/register.html`
```html
<form action="/user/register" method="POST">
    <div class="container">
        <label><b>Firstname</b></label>
        <input type="text" placeholder="Enter Firstname" name="firstname" required>

        <label><b>Username</b></label>
        <input type="text" placeholder="Enter Username" name="username" required>

        <label><b>Password</b></label>
        <input type="password" placeholder="Enter Password" name="password" required>

        <button type="submit">Register</button>
    </div>
</form>
```
> `main.go`
```go
package main

import (
	"time"

	"github.com/kataras/iris/v12/_examples/structuring/login-mvc-single-responsibility-package/user"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

func main() {
	app := iris.New()
	//您获得了完整的调试消息，这在使用MVC时很有用，并且您希望
	//确保您的代码与Iris的MVC体系结构保持一致。
	
	// You got full debug messages, useful when using MVC and you want to make
	// sure that your code is aligned with the Iris' MVC Architecture.
	app.Logger().SetLevel("debug")

	app.RegisterView(iris.HTML("./views", ".html").Layout("shared/layout.html"))

	app.HandleDir("/public", "./public")

	mvc.Configure(app, configureMVC)

	// http://localhost:8080/user/register
	// http://localhost:8080/user/login
	// http://localhost:8080/user/me
	// http://localhost:8080/user/logout
	// http://localhost:8080/user/1
	app.Run(iris.Addr(":8080"), configure)
}

func configureMVC(app *mvc.Application) {
	manager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})

	userApp := app.Party("/user")
	userApp.Register(
		user.NewDataSource(),
		manager.Start,
	)
	userApp.Handle(new(user.Controller))
}

func configure(app *iris.Application) {
	app.Configure(
		iris.WithoutServerError(iris.ErrServerClosed),
	)
}
```