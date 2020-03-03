package routes

import (
	"github.com/kataras/iris/v12/_examples/structuring/bootstrap/bootstrap"
)

//Configure将必要的路由注册到应用程序

// Configure registers the necessary routes to the app.
func Configure(b *bootstrap.Bootstrapper) {
	b.Get("/", GetIndexHandler)
	b.Get("/follower/{id:long}", GetFollowerHandler)
	b.Get("/following/{id:long}", GetFollowingHandler)
	b.Get("/like/{id:long}", GetLikeHandler)
}
