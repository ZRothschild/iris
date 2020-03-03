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
