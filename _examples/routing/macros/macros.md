# `route` 路由宏的使用
## 目录结构
> 主目录`macros`
```html
    —— main.go
```
## 代码示例
> `main.go`

```go
// Package main展示了如何注册自定义参数类型和所属的宏函数。
// Package main shows how you can register a custom parameter type and macro functions that belongs to it.
package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/hero"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	app.Macros().Register("slice", "", false, true, func(paramValue string) (interface{}, bool) {
		return strings.Split(paramValue, "/"), true
	}).RegisterFunc("contains", func(expectedItems []string) func(paramValue []string) bool {
		sort.Strings(expectedItems)
		return func(paramValue []string) bool {
			if len(paramValue) != len(expectedItems) {
				return false
			}

			sort.Strings(paramValue)
			for i := 0; i < len(paramValue); i++ {
				if paramValue[i] != expectedItems[i] {
					return false
				}
			}

			return true
		}
	})

	//为了在MVC控制器的函数输入自变量或依赖注入函数输入自变量中使用新的param类型
	//您必须告诉Iris它是什么类型，参数的ValueRaw是相同的类型
	//如上定义的func(paramValue string) (interface{}, bool)一样。
	//现在，新值及其类型（从字符串到新的自定义类型）仅存储一次，
	//对于此类简单情况，您无需进行任何转换

	// In order to use your new param type inside MVC controller's function input argument or a hero function input argument
	// you have to tell the Iris what type it is, the `ValueRaw` of the parameter is the same type
	// as you defined it above with the func(paramValue string) (interface{}, bool).
	// The new value and its type(from string to your new custom type) it is stored only once now,
	// you don't have to do any conversions for simple cases like this.
	context.ParamResolvers[reflect.TypeOf([]string{})] = func(paramIndex int) interface{} {
		return func(ctx iris.Context) []string {
			//如果要检索默认情况下不支持的值类型的参数，例如ctx.Params().GetInt
			//然后，您可以使用`GetEntry`或`GetEntryAt`并将其强制`ValueRaw`转换为所需的类型。
			//类型应与宏函数（Macros＃Register的最后一个参数）的返回值相同。

			// When you want to retrieve a parameter with a value type that it is not supported by-default, such as ctx.Params().GetInt
			// then you can use the `GetEntry` or `GetEntryAt` and cast its underline `ValueRaw` to the desired type.
			// The type should be the same as the macro's evaluator function (last argument on the Macros#Register) return value.
			return ctx.Params().GetEntryAt(paramIndex).ValueRaw.([]string)
		}
	}

	/*
		http://localhost:8080/test_slice_hero/myvaluei1/myavlue2 ->
		myparam的值（尾随路径参数类型）为：[]string{"myvalue1", "myavlue2"}

		myparam's value (a trailing path parameter type) is: []string{"myvalue1", "myavlue2"}
	*/
	app.Get("/test_slice_hero/{myparam:slice}", hero.Handler(func(myparam []string) string {
		return fmt.Sprintf("myparam's value (a trailing path parameter type) is: %#v\n", myparam)
	}))

	/*
		http://localhost:8080/test_slice_contains/notcontains1/value2 ->
		(404) Not Found

		http://localhost:8080/test_slice_contains/value1/value2 ->
		myparam的值（尾随路径参数类型）为：[]string{"value1", "value2"}

		myparam's value (a trailing path parameter type) is: []string{"value1", "value2"}
	*/
	app.Get("/test_slice_contains/{myparam:slice contains([value1,value2])}", func(ctx iris.Context) {
		//如果没有内置函数可用于以所需类型检索值，例如ctx.Params().GetInt
		//然后，您可以使用`GetEntry.ValueRaw`获取实际值，该值由上面的宏设置。
		
		// When it is not a builtin function available to retrieve your value with the type you want, such as ctx.Params().GetInt
		// then you can use the `GetEntry.ValueRaw` to get the real value, which is set-ed by your macro above.
		myparam := ctx.Params().GetEntry("myparam").ValueRaw.([]string)
		ctx.Writef("myparam's value (a trailing path parameter type) is: %#v\n", myparam)
	})

	app.Run(iris.Addr(":8080"))
}
```