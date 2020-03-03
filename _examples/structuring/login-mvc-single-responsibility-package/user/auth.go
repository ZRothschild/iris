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
