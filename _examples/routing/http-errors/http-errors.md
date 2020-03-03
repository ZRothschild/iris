# `route` 自定义HTTP错误处理
## 目录结构
> 主目录`http-errors`
```html
    —— main.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	//捕获特定的错误代码
	// Catch a specific error code.
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.HTML("Message: <b>" + ctx.Values().GetString("message") + "</b>")
	})
	//捕获所有错误代码[app.OnAnyErrorCode...]
	// Catch all error codes [app.OnAnyErrorCode...]
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML(`Click <a href="/my500">here</a> to pretend an HTTP error`)
	})

	app.Get("/my500", func(ctx iris.Context) {
		ctx.Values().Set("message", "this is the error message")
		ctx.StatusCode(500)
	})

	app.Get("/u/{firstname:alphabetical}", func(ctx iris.Context) {
		ctx.Writef("Hello %s", ctx.Params().Get("firstname"))
	})

	// Read more at: https://github.com/kataras/iris/issues/1335
	app.Get("/product-problem", problemExample)

	app.Get("/product-error", func(ctx iris.Context) {
		ctx.Writef("explain the error")
	})

	// http://localhost:8080
	// http://localhost:8080/my500
	// http://localhost:8080/u/gerasimos
	// http://localhost:8080/product-problem
	app.Run(iris.Addr(":8080"))
}

func newProductProblem(productName, detail string) iris.Problem {
	return iris.NewProblem().
		// URI类型，如果是相对类型，它将自动转换为绝对类型
		// The type URI, if relative it automatically convert to absolute.
		Type("/product-error").
		//标题，如果为空，则从状态代码获取
		// The title, if empty then it gets it from the status code.
		Title("Product validation problem").
		//任何可选的详细信息
		// Any optional details.
		Detail(detail).
		//状态错误代码，必填
		// The status error code, required.
		Status(iris.StatusBadRequest).
		//任何自定义键值对
		// Any custom key-value pair.
		Key("productName", productName)
	//问题的可选原因，问题链。
	// Optional cause of the problem, chain of Problems.
	// Cause(iris.NewProblem().Type("/error").Title("cause of the problem").Status(400))
}

func problemExample(ctx iris.Context) {
	/*
		p := iris.NewProblem().
			Type("/validation-error").
			Title("Your request parameters didn't validate").
			Detail("Optional details about the error.").
			Status(iris.StatusBadRequest).
		 	Key("customField1", customValue1)
		 	Key("customField2", customValue2)
		ctx.Problem(p)

		// OR
		ctx.Problem(iris.Problem{
			"type":   "/validation-error",
			"title":  "Your request parameters didn't validate",
			"detail": "Optional details about the error.",
			"status": iris.StatusBadRequest,
		 	"customField1": customValue1,
		 	"customField2": customValue2,
		})

		// OR
	*/
	//响应类似JSON，但缩进为"  "，并且
	//内容类型为"application/problem+json"

	// Response like JSON but with indent of "  " and
	// content type of "application/problem+json"
	ctx.Problem(newProductProblem("product name", "problem error details"), iris.ProblemOptions{
		//可选的JSON渲染器设置。
		// Optional JSON renderer settings.
		JSON: iris.JSON{
			Indent: "  ",
		},
		// 要么
		//呈现为XML：
		//
		// RenderXML：true，
		// XML：      iris.XML{Indent: "  "},
		//和ctx.StatusCode(200)，以用户身份在浏览器上查看结果。
		//
		//以下`RetryAfter`字段设置"Retry-After"响应标头。
		//
		//可以接受：
		// HTTP-Date  time.Time
		// time.Duration，int64，float64，int 的秒类型
		//或字符串日期类型 和 Duration
		// 例子：
		// time.Now().Add(5 * time.Minute),
		// 300 * time.Second,
		// "5m",

		// OR
		// Render as XML:
		//
		// RenderXML: true,
		// XML:       iris.XML{Indent: "  "},
		// and ctx.StatusCode(200) to see the result on browser as a user.
		//
		// The below `RetryAfter` field sets the "Retry-After" response header.
		//
		// Can accept:
		// time.Time for HTTP-Date,
		// time.Duration, int64, float64, int for seconds
		// or string for date or duration.
		// Examples:
		// time.Now().Add(5 * time.Minute),
		// 300 * time.Second,
		// "5m",
		//
		RetryAfter: 300,
		//一个可以动态设置的函数（如果已指定）
		//根据请求重试。 对于ProblemOptions可重用性很有用。
		//覆盖RetryAfter字段。
		
		// A function that, if specified, can dynamically set
		// retry-after based on the request. Useful for ProblemOptions reusability.
		// Overrides the RetryAfter field.
		//
		// RetryAfterFunc: func(iris.Context) interface{} { [...] }
	})
}

```