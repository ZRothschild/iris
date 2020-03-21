# OAuth2.0认证授权
## 认证授权过程
- 服务提供方(provider)，用户使用服务提供方来存储受保护的资源，如照片，视频，联系人列表
- 用户，存放在服务提供方的受保护的资源的拥有者
- 客户端，要访问服务提供方资源的第三方应用，通常是网站，如提供照片打印服务的网站。在认证过程之前，客户端要向服务提供者申请客户端标识

**使用OAuth进行认证和授权的过程如下所示:**
> 用户想操作存放在服务提供方的资源。
> 用户登录客户端向服务提供方请求一个临时令牌。
> 服务提供方验证客户端的身份后，授予一个临时令牌。
> 客户端获得临时令牌后，将用户引导至服务提供方的授权页面请求用户授权。在这个过程中将临时令牌和客户端的回调连接发送给服务提供方。
> 用户在服务提供方的网页上输入用户名和密码，然后授权该客户端访问所请求的资源。
> 授权成功后，服务提供方引导用户返回客户端的网页。
> 客户端根据临时令牌从服务提供方那里获取访问令牌。
> 服务提供方根据临时令牌和用户的授权情况授予客户端访问令牌。
> 客户端使用获取的访问令牌访问存放在服务提供方上的受保护的资源。
## 目录结构
> 主目录`oauth2`
```html
    —— main.go
    —— templates
        —— index.html
        —— user.html
```
## 代码示例
> `main.go`
```golang
package main

//任何OAuth2（甚至是纯golang/x/net/oauth2）包都可以与iris一起使用，但在此示例中，我们将看到markbates的goth：
// $ go get github.com/markbates/goth/...
//
//此OAuth2示例适用于会话，因此我们需要附加一个会话管理器。 可选：为了获得更安全的会话值，
//开发人员可以使用任何第三方程序包添加自定义cookie编码器/解码器。
//在此示例中，我们将使用gorilla的securecookie：
//
// $ go get github.com/gorilla/securecookie
//可以在"sessions/securecookie"示例文件夹中找到securecookie的示例。

//注意：
//整个示例由markbates/goth/example/main.go转换
//它已经使用我自己的TWITTER应用程序进行了测试，并且即使对于localhost也可以使用
//我想其他一切都会按预期进行，在我编写该示例时，goth库社区报告的所有bug都已修复，玩得开心！

// Any OAuth2 (even the pure golang/x/net/oauth2) package
// can be used with iris but at this example we will see the markbates' goth:
//
// $ go get github.com/markbates/goth/...
//
// This OAuth2 example works with sessions, so we will need
// to attach a session manager.
// Optionally: for even more secure session values,
// developers can use any third-party package to add a custom cookie encoder/decoder.
// At this example we will use the gorilla's securecookie:
//
// $ go get github.com/gorilla/securecookie
// Example of securecookie can be found at "sessions/securecookie" example folder.

// Notes:
// The whole example is converted by markbates/goth/example/main.go.
// It's tested with my own TWITTER application and it worked, even for localhost.
// I guess that everything else works as expected, all bugs reported by goth library's community
// are fixed in the time I wrote that example, have fun!
import (
	"errors"
	"os"
	"sort"

	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/sessions"
	//可选，用于seesion的encoder/decoder
	"github.com/gorilla/securecookie" // optionally, used for session's encoder/decoder

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/dailymotion"
	"github.com/markbates/goth/providers/deezer"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/fitbit"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/meetup"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/stripe"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/xero"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
)

var sessionsManager *sessions.Sessions

func init() {
	//附加session管理器

	// attach a session manager
	cookieName := "mycustomsessionid"
	// AES仅支持16、24或32字节的密钥大小
	// 您要么需要提供准确的字节数，要么从键入的内容中得出密钥

	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	hashKey := []byte("the-big-and-secret-fash-key-here")
	blockKey := []byte("lot-secret-of-characters-big-too")
	secureCookie := securecookie.New(hashKey, blockKey)

	sessionsManager = sessions.New(sessions.Config{
		Cookie: cookieName,
		Encode: secureCookie.Encode,
		Decode: secureCookie.Decode,
	})
}
//这些是您可以使用的一些函数帮助器

// These are some function helpers that you may use if you want

// GetProviderName是用于获取给定请求的提供程序名称的函数。默认情况下，此提供程序是从URL查询字符串中获取的
// 如果以其他方式提供它，请将您自己的函数分配给该变量，该变量返回您请求的提供者名称

// GetProviderName is a function used to get the name of a provider
// for a given request. By default, this provider is fetched from
// the URL query string. If you provide it in a different way,
// assign your own function to this variable that returns the provider
// name for your request.
var GetProviderName = func(ctx iris.Context) (string, error) {
	//尝试从url参数"provider"中获取

	// try to get it from the url param "provider"
	if p := ctx.URLParam("provider"); p != "" {
		return p, nil
	}
	//尝试从url PATH参数“{provider}或：provider或{provider：string}或{provider：alphabetical}”获取它

	// try to get it from the url PATH parameter "{provider} or :provider or {provider:string} or {provider:alphabetical}"
	if p := ctx.Params().Get("provider"); p != "" {
		return p, nil
	}
	//尝试从上下文的每个请求存储中获取它

	// try to get it from context's per-request storage
	if p := ctx.Values().GetString("provider"); p != "" {
		return p, nil
	}
	//如果没有找到，则返回一个带有相应错误的空字符串

	// if not found then return an empty string with the corresponding error
	return "", errors.New("you must select a provider")
}

/*
BeginAuthHandler是用于启动身份验证过程的便捷处理程序
它希望能够从查询参数中获取提供程序的名称授权应用名称）
作为“provider”或“：provider”。

BeginAuthHandler会将用户重定向到相应的身份验证端点对于请求的provider。
请参阅https://github.com/markbates/goth/examples/main.go以查看此操作

BeginAuthHandler is a convenience handler for starting the authentication process.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".

BeginAuthHandler will redirect the user to the appropriate authentication end-point
for the requested provider.

See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
func BeginAuthHandler(ctx iris.Context) {
	url, err := GetAuthURL(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("%v", err)
		return
	}

	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}

/*
GetAuthURL使用provided的请求启动身份验证过程。
它将返回一个应该用于向用户发送的URL。

它希望能够从查询参数中获取provider的名称
作为"provider" 或 ":provider"或来自“provider” key的上下文值。

我建议使用BeginAuthHandler而不是执行所有这些步骤，但那完全取决于你

GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.

It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider" or from the context's value of "provider" key.

I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func GetAuthURL(ctx iris.Context) (string, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	session := sessionsManager.Start(ctx)
	session.Set(providerName, sess.Marshal())
	return url, nil
}

// SetState设置与给定请求关联的状态字符串
// 如果没有状态字符串与请求关联，则将生成一个
// 此状态将发送到提供程序，并且可以在回调期间进行检索

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(ctx iris.Context) string {
	state := ctx.URLParam("state")
	if len(state) > 0 {
		return state
	}

	return "state"
}

// GetState获取提供者在回调过程中返回的状态
// 此状态用于防止CSRF攻击，请参阅
// http://tools.ietf.org/html/rfc6749#section-10.12

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx iris.Context) string {
	return ctx.URLParam("state")
}

/*
CompleteUserAuth会按照提示说。 它完成了身份验证过程，并从提供者那里获取了有关用户的所有基本信息

它期望能够从查询参数中以"provider"或":provider"的形式获取提供者的名称

请参阅https://github.com/markbates/goth/examples/main.go以查看实际效果

CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.

It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".

See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
var CompleteUserAuth = func(ctx iris.Context) (goth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}
	session := sessionsManager.Start(ctx)
	value := session.GetString(providerName)
	if value == "" {
		return goth.User{}, errors.New("session value for " + providerName + " not found")
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		//可以找到具有现有会话数据的用户

		// user can be found with existing session data
		return user, err
	}
	//获取新令牌并重试获取

	// get new token and retry fetch
	_, err = sess.Authorize(provider, ctx.Request().URL.Query())
	if err != nil {
		return goth.User{}, err
	}

	session.Set(providerName, sess.Marshal())
	return provider.FetchUser(sess)
}

//注销会使用户会话无效

// Logout invalidates a user session.
func Logout(ctx iris.Context) error {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return err
	}
	session := sessionsManager.Start(ctx)
	session.Delete(providerName)
	return nil
}

//“一些函数助手”的结尾

// End of the "some function helpers".

func main() {
	goth.UseProviders(
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),
		//如果您想使用身份验证而不是在Twitter提供程序中进行授权，请改用它

		// If you'd like to use authenticate instead of authorize in Twitter provider, use this instead.
		// twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),

		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"),
		fitbit.New(os.Getenv("FITBIT_KEY"), os.Getenv("FITBIT_SECRET"), "http://localhost:3000/auth/fitbit/callback"),
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), "http://localhost:3000/auth/gplus/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
		spotify.New(os.Getenv("SPOTIFY_KEY"), os.Getenv("SPOTIFY_SECRET"), "http://localhost:3000/auth/spotify/callback"),
		linkedin.New(os.Getenv("LINKEDIN_KEY"), os.Getenv("LINKEDIN_SECRET"), "http://localhost:3000/auth/linkedin/callback"),
		lastfm.New(os.Getenv("LASTFM_KEY"), os.Getenv("LASTFM_SECRET"), "http://localhost:3000/auth/lastfm/callback"),
		twitch.New(os.Getenv("TWITCH_KEY"), os.Getenv("TWITCH_SECRET"), "http://localhost:3000/auth/twitch/callback"),
		dropbox.New(os.Getenv("DROPBOX_KEY"), os.Getenv("DROPBOX_SECRET"), "http://localhost:3000/auth/dropbox/callback"),
		digitalocean.New(os.Getenv("DIGITALOCEAN_KEY"), os.Getenv("DIGITALOCEAN_SECRET"), "http://localhost:3000/auth/digitalocean/callback", "read"),
		bitbucket.New(os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"), "http://localhost:3000/auth/bitbucket/callback"),
		instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), "http://localhost:3000/auth/instagram/callback"),
		intercom.New(os.Getenv("INTERCOM_KEY"), os.Getenv("INTERCOM_SECRET"), "http://localhost:3000/auth/intercom/callback"),
		box.New(os.Getenv("BOX_KEY"), os.Getenv("BOX_SECRET"), "http://localhost:3000/auth/box/callback"),
		salesforce.New(os.Getenv("SALESFORCE_KEY"), os.Getenv("SALESFORCE_SECRET"), "http://localhost:3000/auth/salesforce/callback"),
		amazon.New(os.Getenv("AMAZON_KEY"), os.Getenv("AMAZON_SECRET"), "http://localhost:3000/auth/amazon/callback"),
		yammer.New(os.Getenv("YAMMER_KEY"), os.Getenv("YAMMER_SECRET"), "http://localhost:3000/auth/yammer/callback"),
		onedrive.New(os.Getenv("ONEDRIVE_KEY"), os.Getenv("ONEDRIVE_SECRET"), "http://localhost:3000/auth/onedrive/callback"),

		//通过代理将localhost.com指向 http://localhost:3000/auth/yahoo/callback，因为yahoo不允许将自定义端口放入重定向uri

		// Pointed localhost.com to http://localhost:3000/auth/yahoo/callback through proxy as yahoo
		// does not allow to put custom ports in redirection uri
		yahoo.New(os.Getenv("YAHOO_KEY"), os.Getenv("YAHOO_SECRET"), "http://localhost.com"),
		slack.New(os.Getenv("SLACK_KEY"), os.Getenv("SLACK_SECRET"), "http://localhost:3000/auth/slack/callback"),
		stripe.New(os.Getenv("STRIPE_KEY"), os.Getenv("STRIPE_SECRET"), "http://localhost:3000/auth/stripe/callback"),
		wepay.New(os.Getenv("WEPAY_KEY"), os.Getenv("WEPAY_SECRET"), "http://localhost:3000/auth/wepay/callback", "view_user"),
		//默认情况下，将使用贝宝生产授权url，请将PAYPAL_ENV=sandbox设置为环境变量，以便在沙盒环境中进行测试

		// By default paypal production auth urls will be used, please set PAYPAL_ENV=sandbox as environment variable for testing
		// in sandbox environment
		paypal.New(os.Getenv("PAYPAL_KEY"), os.Getenv("PAYPAL_SECRET"), "http://localhost:3000/auth/paypal/callback"),
		steam.New(os.Getenv("STEAM_KEY"), "http://localhost:3000/auth/steam/callback"),
		heroku.New(os.Getenv("HEROKU_KEY"), os.Getenv("HEROKU_SECRET"), "http://localhost:3000/auth/heroku/callback"),
		uber.New(os.Getenv("UBER_KEY"), os.Getenv("UBER_SECRET"), "http://localhost:3000/auth/uber/callback"),
		soundcloud.New(os.Getenv("SOUNDCLOUD_KEY"), os.Getenv("SOUNDCLOUD_SECRET"), "http://localhost:3000/auth/soundcloud/callback"),
		gitlab.New(os.Getenv("GITLAB_KEY"), os.Getenv("GITLAB_SECRET"), "http://localhost:3000/auth/gitlab/callback"),
		dailymotion.New(os.Getenv("DAILYMOTION_KEY"), os.Getenv("DAILYMOTION_SECRET"), "http://localhost:3000/auth/dailymotion/callback", "email"),
		deezer.New(os.Getenv("DEEZER_KEY"), os.Getenv("DEEZER_SECRET"), "http://localhost:3000/auth/deezer/callback", "email"),
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), "http://localhost:3000/auth/discord/callback", discord.ScopeIdentify, discord.ScopeEmail),
		meetup.New(os.Getenv("MEETUP_KEY"), os.Getenv("MEETUP_SECRET"), "http://localhost:3000/auth/meetup/callback"),

		// Auth0为每个客户分配域，必须提供一个域才能使auth0正常工作

		// Auth0 allocates domain per customer, a domain must be provided for auth0 to work
		auth0.New(os.Getenv("AUTH0_KEY"), os.Getenv("AUTH0_SECRET"), "http://localhost:3000/auth/auth0/callback", os.Getenv("AUTH0_DOMAIN")),
		xero.New(os.Getenv("XERO_KEY"), os.Getenv("XERO_SECRET"), "http://localhost:3000/auth/xero/callback"),
	)

	// OpenID Connect基于OpenID Connect自动发现URL（https://openid.net/specs/openid-connect-discovery-1_0-17.html）
	// 因为OpenID Connect提供程序会在New()中自行对其进行初始化,可以返回应该处理或忽略的错误，暂时忽略该错误

	// OpenID Connect is based on OpenID Connect Auto Discovery URL (https://openid.net/specs/openid-connect-discovery-1_0-17.html)
	// because the OpenID Connect provider initialize it self in the New(), it can return an error which should be handled or ignored
	// ignore the error for now
	openidConnect, _ := openidConnect.New(os.Getenv("OPENID_CONNECT_KEY"), os.Getenv("OPENID_CONNECT_SECRET"), "http://localhost:3000/auth/openid-connect/callback", os.Getenv("OPENID_CONNECT_DISCOVERY_URL"))
	if openidConnect != nil {
		goth.UseProviders(openidConnect)
	}

	m := make(map[string]string)
	m["amazon"] = "Amazon"
	m["bitbucket"] = "Bitbucket"
	m["box"] = "Box"
	m["dailymotion"] = "Dailymotion"
	m["deezer"] = "Deezer"
	m["digitalocean"] = "Digital Ocean"
	m["discord"] = "Discord"
	m["dropbox"] = "Dropbox"
	m["facebook"] = "Facebook"
	m["fitbit"] = "Fitbit"
	m["github"] = "Github"
	m["gitlab"] = "Gitlab"
	m["soundcloud"] = "SoundCloud"
	m["spotify"] = "Spotify"
	m["steam"] = "Steam"
	m["stripe"] = "Stripe"
	m["twitch"] = "Twitch"
	m["uber"] = "Uber"
	m["wepay"] = "Wepay"
	m["yahoo"] = "Yahoo"
	m["yammer"] = "Yammer"
	m["gplus"] = "Google Plus"
	m["heroku"] = "Heroku"
	m["instagram"] = "Instagram"
	m["intercom"] = "Intercom"
	m["lastfm"] = "Last FM"
	m["linkedin"] = "Linkedin"
	m["onedrive"] = "Onedrive"
	m["paypal"] = "Paypal"
	m["twitter"] = "Twitter"
	m["salesforce"] = "Salesforce"
	m["slack"] = "Slack"
	m["meetup"] = "Meetup.com"
	m["auth0"] = "Auth0"
	m["openid-connect"] = "OpenID Connect"
	m["xero"] = "Xero"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

	//创建我们的应用，
	//设置一个视图
	//设置会话
	//并为展示柜设置路由器

	// create our app,
	// set a view
	// set sessions
	// and setup the router for the showcase
	app := iris.New()

	//附加并构建我们的模板

	// attach and build our templates
	app.RegisterView(iris.HTML("./templates", ".html"))

	//路由器的启动

	// start of the router

	app.Get("/auth/{provider}/callback", func(ctx iris.Context) {
		user, err := CompleteUserAuth(ctx)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Writef("%v", err)
			return
		}
		ctx.ViewData("", user)
		if err := ctx.View("user.html"); err != nil {
			ctx.Writef("%v", err)
		}
	})

	app.Get("/logout/{provider}", func(ctx iris.Context) {
		Logout(ctx)
		ctx.Redirect("/", iris.StatusTemporaryRedirect)
	})

	app.Get("/auth/{provider}", func(ctx iris.Context) {
		//尝试在不重新认证的情况下获取用户
		
		// try to get the user without re-authenticating
		if gothUser, err := CompleteUserAuth(ctx); err == nil {
			ctx.ViewData("", gothUser)
			if err := ctx.View("user.html"); err != nil {
				ctx.Writef("%v", err)
			}
		} else {
			BeginAuthHandler(ctx)
		}
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("", providerIndex)

		if err := ctx.View("index.html"); err != nil {
			ctx.Writef("%v", err)
		}
	})

	// http://localhost:3000
	app.Run(iris.Addr("localhost:3000"))
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}
```
> `index.html`
```html
{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}
```
> `user.html`
```html
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
```
### 提示
1. 去Github设置第三方授权登录，设置回调路径，把key,secret填入上面代码指定位置
2. 执行程序选择github 即可看到想要的效果