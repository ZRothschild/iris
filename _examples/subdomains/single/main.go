//打包主注册静态子域，如果使用Windows，请检查./hosts
// Package main register static subdomains, simple as parties, check ./hosts if you use windows
package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	//子域可与所有可用路由器一起使用，就像其他功能一样

	// subdomains works with all available routers, like other features too.

	//没有顺序，您也可以在末尾注册子域

	// no order, you can register subdomains at the end also.
	admin := app.Party("admin.")
	{
		// admin.mydomain.com
		admin.Get("/", func(c iris.Context) {
			c.Writef("INDEX FROM admin.mydomain.com")
		})
		// admin.mydomain.com/hey
		admin.Get("/hey", func(c iris.Context) {
			c.Writef("HEY FROM admin.mydomain.com/hey")
		})
		// admin.mydomain.com/hey2
		admin.Get("/hey2", func(c iris.Context) {
			c.Writef("HEY SECOND FROM admin.mydomain.com/hey")
		})
	}

	// mydomain.com/
	app.Get("/", func(c iris.Context) {
		c.Writef("INDEX FROM no-subdomain hey")
	})

	// mydomain.com/hey
	app.Get("/hey", func(c iris.Context) {
		c.Writef("HEY FROM no-subdomain hey")
	})

	// http://admin.mydomain.com
	// http://admin.mydomain.com/hey
	// http://admin.mydomain.com/hey2
	// http://mydomain.com
	// http://mydomain.com/hey

	//对于初学者：查看../hosts文件
	app.Run(iris.Addr("mydomain.com:80")) // for beginners: look ../hosts file
}
