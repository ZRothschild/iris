# IRIS路由动态匹配
## 目录结构
> 主目录`dynamic-path`
```html
    —— main.go
    —— root-wildcard
        —— main.go
```
## 代码示例
> `main.go`

```golang
package main

import (
	"regexp"

	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	//在上一个示例"routing/basic"中，
	//我们已经看到了静态路由，路由组，子域，通配符子域，参数化路径的一个小示例
	//具有一个已知的参数和自定义http错误，现在是时候查看通配符参数和宏了。

	//Iris，如net/http std包注册路由的处理程序
	//通过处理程序，iris的处理程序类型仅仅是func(ctx iris.Context)
	//上下文来自github.com/kataras/iris/context。
	//
	// Iris具有您遇到过的最简单，最强大的路由过程。
	//
	// 同时，
	//iris有它自己的交互器（是像编程语言一样）
	//用于路由的路径语法及其动态路径参数的解析和评估，
	//我们称它们为“宏”是捷径。
	// 怎么样？它计算其需求，如果不需要任何特殊的正则表达式，那么它只是
	//使用低级下划线路径语法注册路由，
	//否则，它将预编译regexp并添加必要的中间件。

	// At the previous example "routing/basic",
	// we've seen static routes, group of routes, subdomains, wildcard subdomains, a small example of parameterized path
	// with a single known paramete and custom http errors, now it's time to see wildcard parameters and macros.

	// Iris, like net/http std package registers route's handlers
	// by a Handler, the iris' type of handler is just a func(ctx iris.Context)
	// where context comes from github.com/kataras/iris/context.
	//
	// Iris has the easiest and the most powerful routing process you have ever meet.
	//
	// At the same time,
	// Iris has its own interpeter(yes like a programming language)
	// for route's path syntax and their dynamic path parameters parsing and evaluation,
	// We call them "macros" for shortcut.
	// How? It calculates its needs and if not any special regexp needed then it just
	// registers the route with the low-level underline  path syntax,
	// otherwise it pre-compiles the regexp and adds the necessary middleware(s).


	//参数的标准宏类型：
	// Standard macro types for parameters:
	// +------------------------+
	// | {param:string}         |
	// +------------------------+
	// string type
	// anything (single path segmnent)
	//
	// +-------------------------------+
	// | {param:int}                   |
	// +-------------------------------+
	// int type
	// -9223372036854775808 to 9223372036854775807 (x64) or -2147483648 to 2147483647 (x32), depends on the host arch
	//
	// +------------------------+
	// | {param:int8}           |
	// +------------------------+
	// int8 type
	// -128 to 127
	//
	// +------------------------+
	// | {param:int16}          |
	// +------------------------+
	// int16 type
	// -32768 to 32767
	//
	// +------------------------+
	// | {param:int32}          |
	// +------------------------+
	// int32 type
	// -2147483648 to 2147483647
	//
	// +------------------------+
	// | {param:int64}          |
	// +------------------------+
	// int64 type
	// -9223372036854775808 to 9223372036854775807
	//
	// +------------------------+
	// | {param:uint}           |
	// +------------------------+
	// uint type
	// 0 to 18446744073709551615 (x64) or 0 to 4294967295 (x32)
	//
	// +------------------------+
	// | {param:uint8}          |
	// +------------------------+
	// uint8 type
	// 0 to 255
	//
	// +------------------------+
	// | {param:uint16}         |
	// +------------------------+
	// uint16 type
	// 0 to 65535
	//
	// +------------------------+
	// | {param:uint32}          |
	// +------------------------+
	// uint32 type
	// 0 to 4294967295
	//
	// +------------------------+
	// | {param:uint64}         |
	// +------------------------+
	// uint64 type
	// 0 to 18446744073709551615
	//
	// +---------------------------------+
	// | {param:bool} or {param:boolean} |
	// +---------------------------------+
	// bool type
	// only "1" or "t" or "T" or "TRUE" or "true" or "True"
	// or "0" or "f" or "F" or "FALSE" or "false" or "False"
	//
	// +------------------------+
	// | {param:alphabetical}   |
	// +------------------------+
	//字母/字母类型
	//仅字母（大写或小写）
	// alphabetical/letter type
	// letters only (upper or lowercase)
	//
	// +------------------------+
	// | {param:file}           |
	// +------------------------+
	// file type
	// letters (upper or lowercase)
	// numbers (0-9)
	// underscore (_)
	// dash (-)
	// point (.)
	// no spaces ! or other character
	//
	// +------------------------+
	// | {param:path}           |
	// +------------------------+

	//路径类型
	//任何东西都应该是最后一部分，可以是多个路径段，
	//即："/test/{param:path}"并请求："/test/path1/path2/path3" , ctx.Params().Get("param") == "path1/path2/path3"
	//
	//如果缺少类型，则参数的类型默认为字符串，因此
	// {param} == {param：string}。
	//
	//如果找不到该类型的函数，则使用`string`宏类型的函数。
	//
	//
	//除了iris提供基本类型和一些默认的“宏功能”
	//您也可以注册自己的！
	//
	//注册一个命名路径参数函数：
	// app.Macros().Number.RegisterFunc("min", func(argument int) func(paramValue string) bool {
	//  [...]
	//  return true/false -> true表示有效
	// })
	//
	//在func(argument ...)中，您可以具有任何标准类型，它将在服务器启动之前进行验证
	//因此，这里不必关心性能，它在服务时运行的唯一内容就是返回的func(paramValue string) bool.
	//
	//{param:string equal(iris)} , "iris"将是此处的参数：
	// app.Macros().String.RegisterFunc("equal", func(argument string) func(paramValue string) bool {
	// 	return func(paramValue string) bool { return argument == paramValue }
	// })

	// path type
	// anything, should be the last part, can be more than one path segment,
	// i.e: "/test/{param:path}" and request: "/test/path1/path2/path3" , ctx.Params().Get("param") == "path1/path2/path3"
	//
	// if type is missing then parameter's type is defaulted to string, so
	// {param} == {param:string}.
	//
	// If a function not found on that type then the `string` macro type's functions are being used.
	//
	//
	// Besides the fact that iris provides the basic types and some default "macro funcs"
	// you are able to register your own too!.
	//
	// Register a named path parameter function:
	// app.Macros().Number.RegisterFunc("min", func(argument int) func(paramValue string) bool {
	//  [...]
	//  return true/false -> true means valid.
	// })
	//
	// at the func(argument ...) you can have any standard type, it will be validated before the server starts
	// so don't care about performance here, the only thing it runs at serve time is the returning func(paramValue string) bool.
	//
	// {param:string equal(iris)} , "iris" will be the argument here:
	// app.Macros().String.RegisterFunc("equal", func(argument string) func(paramValue string) bool {
	// 	return func(paramValue string) bool { return argument == paramValue }
	// })

	//您可以使用"string"类型，该类型对于单个路径参数（可以是任意值）有效

	// you can use the "string" type which is valid for a single path parameter that can be anything.
	app.Get("/username/{name}", func(ctx iris.Context) {
		ctx.Writef("Hello %s", ctx.Params().Get("name"))
	}) // type is missing = {name:string}

	//让我们注册连接到uint64宏类型的第一个宏。
	//"min" =  =函数
	//"minValue" =函数的参数
	// func(uint64) bool =我们的func的求值器，它在投放时执行
	//用户请求满足min(...)数方法且满足：uint64参数类型的路径。

	// Let's register our first macro attached to uint64 macro type.
	// "min" = the function
	// "minValue" = the argument of the function
	// func(uint64) bool = our func's evaluator, this executes in serve time when
	// a user requests a path which contains the :uint64 macro parameter type with the min(...) macro parameter function.
	app.Macros().Get("uint64").RegisterFunc("min", func(minValue uint64) func(uint64) bool {
		//"paramValue"的类型应该与内部宏的函数参数的类型一样，在本例中为“ uint64”。
		// type of "paramValue" should match the type of the internal macro's evaluator function, which in this case is "uint64".
		return func(paramValue uint64) bool {
			return paramValue >= minValue
		}
	})

	// http://localhost:8080/profile/id>=20
	//即使在/profile/0, /profile/blabla, /profile/-1上找到路由，这也会抛出404
	//宏参数函数当然是可选的。
	// this will throw 404 even if it's found as route on : /profile/0, /profile/blabla, /profile/-1
	// macro parameter functions are optional of course.
	app.Get("/profile/{id:uint64 min(20)}", func(ctx iris.Context) {
		//第二个参数是错误，但由于使用宏，它始终为nil，
		//验证已经发生。

		// second parameter is the error but it will always nil because we use macros,
		// the validaton already happened.
		id := ctx.Params().GetUint64Default("id", 0)
		ctx.Writef("Hello id: %d", id)
	})

	//更改每个路由的宏处理程序的错误代码：
	// to change the error code per route's macro evaluator:
	app.Get("/profile/{id:uint64 min(1)}/friends/{friendid:uint64 min(1) else 504}", func(ctx iris.Context) {
		id := ctx.Params().GetUint64Default("id", 0)
		friendid := ctx.Params().GetUint64Default("friendid", 0)
		ctx.Writef("Hello id: %d looking for friend id: ", id, friendid)
		//如果路由未通过所有宏处理程序，则会抛出504错误代码而不是404错误代码。
	}) // this will throw e 504 error code instead of 404 if all route's macros not passed.

	// :uint8 0 to 255.
	app.Get("/ages/{age:uint8 else 400}", func(ctx iris.Context) {
		age, _ := ctx.Params().GetUint8("age")
		ctx.Writef("age selected: %d", age)
	})

	//使用自定义正则表达式或任何自定义逻辑的另一个示例。

	//将自定义的无参数宏函数注册为：string参数类型。

	// Another example using a custom regexp or any custom logic.

	// Register your custom argument-less macro function to the :string param type.
	latLonExpr := "^-?[0-9]{1,3}(?:\\.[0-9]{1,10})?$"
	latLonRegex, err := regexp.Compile(latLonExpr)
	if err != nil {
		panic(err)
	}
	// MatchString是func(string) bool的一种，因此我们按原样使用它

	// MatchString is a type of func(string) bool, so we use it as it is.
	app.Macros().Get("string").RegisterFunc("coordinate", latLonRegex.MatchString)

	app.Get("/coordinates/{lat:string coordinate() else 502}/{lon:string coordinate() else 502}", func(ctx iris.Context) {
		ctx.Writef("Lat: %s | Lon: %s", ctx.Params().Get("lat"), ctx.Params().Get("lon"))
	})

	//
	//另一个是通过使用自定义示例
	// Another one is by using a custom body.
	app.Macros().Get("string").RegisterFunc("range", func(minLength, maxLength int) func(string) bool {
		return func(paramValue string) bool {
			return len(paramValue) >= minLength && len(paramValue) <= maxLength
		}
	})

	app.Get("/limitchar/{name:string range(1,200)}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef(`Hello %s | the name should be between 1 and 200 characters length
		otherwise this handler will not be executed`, name)
	})

	//
	//注册您的自定义宏函数，该宏函数接受字符串[`[...,...]`的一部分。
	// Register your custom macro function which accepts a slice of strings `[...,...]`.
	app.Macros().Get("string").RegisterFunc("has", func(validNames []string) func(string) bool {
		return func(paramValue string) bool {
			for _, validName := range validNames {
				if validName == paramValue {
					return true
				}
			}

			return false
		}
	})

	app.Get("/static_validation/{name:string has([kataras,gerasimos,maropoulos])}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef(`Hello %s | the name should be "kataras" or "gerasimos" or "maropoulos"
		otherwise this handler will not be executed`, name)
	})

	//

	// http://localhost:8080/game/a-zA-Z/level/42
	//请记住，alphabetical仅是小写或大写字母。
	// remember, alphabetical is lowercase or uppercase letters only.
	app.Get("/game/{name:alphabetical}/level/{level:int}", func(ctx iris.Context) {
		ctx.Writef("name: %s | level: %s", ctx.Params().Get("name"), ctx.Params().Get("level"))
	})

	app.Get("/lowercase/static", func(ctx iris.Context) {
		ctx.Writef("static and dynamic paths are not conflicted anymore!")
	})

	//让我们使用一个简单的自定义正则表达式来验证单个路径参数
	//其值仅是小写字母。

	// let's use a trivial custom regexp that validates a single path parameter
	// which its value is only lowercase letters.

	// http://localhost:8080/lowercase/anylowercase
	app.Get("/lowercase/{name:string regexp(^[a-z]+)}", func(ctx iris.Context) {
		ctx.Writef("name should be only lowercase, otherwise this handler will never executed: %s", ctx.Params().Get("name"))
	})

	// http://localhost:8080/single_file/app.js
	app.Get("/single_file/{myfile:file}", func(ctx iris.Context) {
		ctx.Writef("file type validates if the parameter value has a form of a file name, got: %s", ctx.Params().Get("myfile"))
	})

	// http://localhost:8080/myfiles/any/directory/here/
	//这是唯一接受任意数量路径段的宏类型。
	// this is the only macro type that accepts any number of path segments.
	app.Get("/myfiles/{directory:path}", func(ctx iris.Context) {
		ctx.Writef("path type accepts any number of path segments, path after /myfiles/ is: %s", ctx.Params().Get("directory"))
		//对于未经验证的通配符路径（任意数量的路径段），您可以使用：/myfiles/*
	}) // for wildcard path (any number of path segments) without validation you can use:
	// /myfiles/*

	//"{param}"的性能与":param"的性能完全相同
	// "{param}"'s performance is exactly the same of ":param"'s.

	//替代->":param"表示单个路径参数，"*"表示通配符路径参数。
	//注意这些：
	//如果为"/mypath/*"，则参数名称为"*"。
	//如果为"/mypath/{myparam:path}"，则该参数具有两个名称，一个是"*"，另一个是用户定义的"myparam".

	// alternatives -> ":param" for single path parameter and "*" for wildcard path parameter.
	// Note these:
	// if  "/mypath/*" then the parameter name is "*".
	// if  "/mypath/{myparam:path}" then the parameter has two names, one is the "*" and the other is the user-defined "myparam".

	// 警告：
	//路径参数名称应仅包含字母或数字。 不允许使用'_'之类的符号。
	//最后，请勿将`ctx.Params()`与`ctx.Values()`混淆。
	//可以从`ctx.Params()`中检索路径参数的值，
	//可以在处理程序和中间件之间进行通信的上下文的本地存储可以存储到`ctx.Values()`中。
	//
	//在相同的确切路径模式中注册不同的参数类型时，路径参数的名称
	//应该有所不同，例如
	// /path/{name:string}
	// /path/{id:uint}
	
	// WARNING:
	// A path parameter name should contain only alphabetical letters or digits. Symbols like  '_' are NOT allowed.
	// Last, do not confuse `ctx.Params()` with `ctx.Values()`.
	// Path parameter's values can be retrieved from `ctx.Params()`,
	// context's local storage that can be used to communicate between handlers and middleware(s) can be stored to `ctx.Values()`.
	//
	// When registering different parameter types in the same exact path pattern, the path parameter's name
	// should differ e.g.
	// /path/{name:string}
	// /path/{id:uint}
	app.Run(iris.Addr(":8080"))
}
```

> `root-wildcard/main.go`

```golang
package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	//现在可以正常工作了，
	//将处理所有GET请求
	// 除了：
	// /                     -> 因为 app.Get("/", ...)
	// /other/anything/here  -> 因为 app.Get("/other/{paramother:path}", ...)
	// /other2/anything/here -> 因为 app.Get("/other2/{paramothersecond:path}", ...)
	// /other2/static2       -> 因为 app.Get("/other2/static", ...)
	//
	//它与其余路由没有冲突，没有路由性能成本！
	//
	//即 /something/here/that/cannot/be/found/by/other/registered/routes/order/not/matters

	// this works as expected now,
	// will handle all GET requests
	// except:
	// /                     -> because of app.Get("/", ...)
	// /other/anything/here  -> because of app.Get("/other/{paramother:path}", ...)
	// /other2/anything/here -> because of app.Get("/other2/{paramothersecond:path}", ...)
	// /other2/static2        -> because of app.Get("/other2/static", ...)
	//
	// It isn't conflicts with the rest of the routes, without routing performance cost!
	//
	// i.e /something/here/that/cannot/be/found/by/other/registered/routes/order/not/matters
	app.Get("/{p:path}", h)
	// app.Get("/static/{p:path}", staticWildcardH)

	//这只会处理 GET /
	// this will handle only GET /
	app.Get("/", staticPath)
	
	//这将处理所有以"/other/"开头的GET请求
	// this will handle all GET requests starting with "/other/"
	//
	// i.e /other/more/than/one/path/parts
	app.Get("/other/{paramother:path}", other)

	//这将处理所有以"/other2/"开头的GET请求
	// /other2/static 除外（由于下一条静态路由）

	// this will handle all GET requests starting with "/other2/"
	// except /other2/static (because of the next static route)
	//
	// i.e /other2/more/than/one/path/parts
	app.Get("/other2/{paramothersecond:path}", other2)
	//这只会处理GET "/other2/static"
	// this will handle only GET "/other2/static"
	app.Get("/other2/static2", staticPathOther2)

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func h(ctx iris.Context) {
	param := ctx.Params().Get("p")
	ctx.WriteString(param)
}

func staticWildcardH(ctx iris.Context) {
	param := ctx.Params().Get("p")
	ctx.WriteString("from staticWildcardH: param=" + param)
}

func other(ctx iris.Context) {
	param := ctx.Params().Get("paramother")
	ctx.Writef("from other: %s", param)
}

func other2(ctx iris.Context) {
	param := ctx.Params().Get("paramothersecond")
	ctx.Writef("from other2: %s", param)
}

func staticPath(ctx iris.Context) {
	ctx.Writef("from the static path(/): %s", ctx.Path())
}

func staticPathOther2(ctx iris.Context) {
	ctx.Writef("from the static path(/other2/static2): %s", ctx.Path())
}
```