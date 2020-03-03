# 使用Go，Iris和Bolt的URL短链接生成程序服务

Hackernoon文章 : https://medium.com/hackernoon/a-url-shortener-service-using-go-iris-and-bolt-4182f0b00ae7

## 目录结构
> 主目录`url-shortener`
```html
    —— resources
        —— css
            —— style.css
    —— templates
        —— index.html
    —— factory.go
    —— main.go
    —— main_test.go
    —— store.go
```
## 代码示例
> `resources/css/style.css`
```css
body{
    background-color:silver;
}
```
> `templates/index.html`
```html
<html>
<head>
    <meta charset="utf-8">
    <title>Golang URL Shortener</title>
    <link rel="stylesheet" href="/static/css/style.css" />
</head>
<body>
    <h2>Golang URL Shortener</h2>
    <h3>{{ .FORM_RESULT}}</h3>
    <form action="/shorten" method="POST">
        <input type="text" name="url" style="width: 35em;" />
        <input type="submit" value="Shorten!" />
    </form>
    {{ if IsPositive .URL_COUNT }}
        <p>{{ .URL_COUNT }} URLs shortened</p>
    {{ end }}
    <form action="/clear_cache" method="POST">
        <input type="submit" value="Clear DB" />
    </form>
</body>
</html>
```
> `factory.go`
```go
package main

import (
	"net/url"

	"github.com/iris-contrib/go.uuid"
)
//生成类型以生成密钥（短网址）

// Generator the type to generate keys(short urls)
type Generator func() string

// DefaultGenerator是默认的URL生成器

// DefaultGenerator is the defautl url generator
var DefaultGenerator = func() string {
	id, _ := uuid.NewV4()
	return id.String()
}

//工厂负责生成密钥（短网址）

// Factory is responsible to generate keys(short urls)
type Factory struct {
	store     Store
	generator Generator
}
// NewFactory接收一个生成器和一个存储，并返回一个新的URL Factory。

// NewFactory receives a generator and a store and returns a new url Factory.
func NewFactory(generator Generator, store Store) *Factory {
	return &Factory{
		store:     store,
		generator: generator,
	}
}

// Gen生成密钥

// Gen generates the key.
func (f *Factory) Gen(uri string) (key string, err error) {
	//我们不返回已解析的url，因为#hash已转换为uri兼容，并且我们不想一直进行编码/解码，
	// 因此不需要这样做，我们将URL保存为用户期望的值，如果 uri验证已通过

	// we don't return the parsed url because #hash are converted to uri-compatible
	// and we don't want to encode/decode all the time, there is no need for that,
	// we save the url as the user expects if the uri validation passed.
	_, err = url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}

	key = f.generator()
	//确保密钥是唯一的

	// Make sure that the key is unique
	for {
		if v := f.store.Get(key); v == "" {
			break
		}
		key = f.generator()
	}

	return key, nil
}
```
> `main.go`
```go
// Package main展示了如何创建简单的URL Shortener。
//
//文章：https：//medium.com/@kataras/a-url-shortener-service-using-go-iris-and-bolt-4182f0b00ae7
//
// Package main shows how you can create a simple URL Shortener.
//
// Article: https://medium.com/@kataras/a-url-shortener-service-using-go-iris-and-bolt-4182f0b00ae7
//
// $ go get github.com/etcd-io/bbolt
// $ go get github.com/iris-contrib/go.uuid
// $ cd $GOPATH/src/github.com/kataras/iris/_examples/tutorial/url-shortener
// $ go build
// $ ./url-shortener
package main

import (
	"fmt"
	"html/template"

	"github.com/kataras/iris/v12"
)

func main() {
	//为数据库分配一个变量，以便稍后使用

	// assign a variable to the DB so we can use its features later.
	db := NewDB("shortener.db")
	//将该数据库传递给我们的应用程序，以便以后可以使用其他数据库测试整个应用程序。

	// Pass that db to our app, in order to be able to test the whole app with a different database later on.
	app := newApp(db)

	//当服务器关闭时释放"db"连接

	// release the "db" connection when server goes off.
	iris.RegisterOnInterrupt(db.Close)

	app.Run(iris.Addr(":8080"))
}

func newApp(db *DB) *iris.Application {
	app := iris.Default() // or app := iris.New()

	//创建我们的工厂，该工厂是对象创建的管理
	//在我们的Web应用程序和数据库之间

	// create our factory, which is the manager for the object creation.
	// between our web app and the db.
	factory := NewFactory(DefaultGenerator, db)

	//通过HTML std视图引擎为"./templates" 目录的“ * .html”文件提供服务

	// serve the "./templates" directory's "*.html" files with the HTML std view engine.
	tmpl := iris.HTML("./templates", ".html").Reload(true)
	//在此处注册任何模板功能
	//
	//看./templates/index.html#L16

	// register any template func(s) here.
	//
	// Look ./templates/index.html#L16
	tmpl.AddFunc("IsPositive", func(n int) bool {
		if n > 0 {
			return true
		}
		return false
	})

	app.RegisterView(tmpl)
	//提供静态文件（css）

	// Serve static files (css)
	app.HandleDir("/static", "./resources")

	indexHandler := func(ctx iris.Context) {
		ctx.ViewData("URL_COUNT", db.Len())
		ctx.View("index.html")
	}
	app.Get("/", indexHandler)

	//通过在http://localhost:8080/u/dsaoj41u321dsa上使用的键来查找并执行短网址

	// find and execute a short url by its key
	// used on http://localhost:8080/u/dsaoj41u321dsa
	execShortURL := func(ctx iris.Context, key string) {
		if key == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		value := db.Get(key)
		if value == "" {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.Writef("Short URL for key: '%s' not found", key)
			return
		}

		ctx.Redirect(value, iris.StatusTemporaryRedirect)
	}
	app.Get("/u/{shortkey}", func(ctx iris.Context) {
		execShortURL(ctx, ctx.Params().Get("shortkey"))
	})

	app.Get("/u/3861bc4d-ca57-4cbc-9fe4-9e0e2b50fff4", func(ctx iris.Context) {
		fmt.Printf("%s\n","testsssss")
	})

	//app.Get("/u/{shortkey}", func(ctx iris.Context) {
	//	execShortURL(ctx, ctx.Params().Get("shortkey"))
	//})

	app.Post("/shorten", func(ctx iris.Context) {
		formValue := ctx.FormValue("url")
		if formValue == "" {
			ctx.ViewData("FORM_RESULT", "You need to a enter a URL")
			ctx.StatusCode(iris.StatusLengthRequired)
		} else {
			key, err := factory.Gen(formValue)
			if err != nil {
				ctx.ViewData("FORM_RESULT", "Invalid URL")
				ctx.StatusCode(iris.StatusBadRequest)
			} else {
				if err = db.Set(key, formValue); err != nil {
					ctx.ViewData("FORM_RESULT", "Internal error while saving the URL")
					app.Logger().Infof("while saving URL: " + err.Error())
					ctx.StatusCode(iris.StatusInternalServerError)
				} else {
					ctx.StatusCode(iris.StatusOK)
					shortenURL := "http://" + app.ConfigurationReadOnly().GetVHost() + "/u/" + key
					ctx.ViewData("FORM_RESULT",
						template.HTML("<pre><a target='_new' href='"+shortenURL+"'>"+shortenURL+" </a></pre>"))
				}
			}
		}
		//没有重定向，我们需要FORM_RESULT
		indexHandler(ctx) // no redirect, we need the FORM_RESULT.
	})

	app.Post("/clear_cache", func(ctx iris.Context) {
		db.Clear()
		ctx.Redirect("/")
	})

	return app
}
```
> `main_test.go`
```go
package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/kataras/iris/v12/httptest"
)

// TestURLShortener tests the simple tasks of our url shortener application.
// Note that it's a pure test.
// The rest possible checks is up to you, take it as as an exercise!
func TestURLShortener(t *testing.T) {
	// temp db file
	f, err := ioutil.TempFile("", "shortener")
	if err != nil {
		t.Fatalf("creating temp file for database failed: %v", err)
	}

	db := NewDB(f.Name())
	app := newApp(db)

	e := httptest.New(t, app)
	originalURL := "https://google.com"

	// save
	e.POST("/shorten").
		WithFormField("url", originalURL).Expect().
		Status(httptest.StatusOK).Body().Contains("<pre><a target='_new' href=")

	keys := db.GetByValue(originalURL)
	if got := len(keys); got != 1 {
		t.Fatalf("expected to have 1 key but saved %d short urls", got)
	}

	// get
	e.GET("/u/" + keys[0]).Expect().
		Status(httptest.StatusTemporaryRedirect).Header("Location").Equal(originalURL)

	// save the same again, it should add a new key
	e.POST("/shorten").
		WithFormField("url", originalURL).Expect().
		Status(httptest.StatusOK).Body().Contains("<pre><a target='_new' href=")

	keys2 := db.GetByValue(originalURL)
	if got := len(keys2); got != 1 {
		t.Fatalf("expected to have 1 keys even if we save the same original url but saved %d short urls", got)
	} // the key is the same, so only the first one matters.

	if keys[0] != keys2[0] {
		t.Fatalf("expected keys to be equal if the original url is the same, but got %s = %s ", keys[0], keys2[0])
	}

	// clear db
	e.POST("/clear_cache").Expect().Status(httptest.StatusOK)
	if got := db.Len(); got != 0 {
		t.Fatalf("expected database to have 0 registered objects after /clear_cache but has %d", got)
	}

	// give it some time to release the db connection
	db.Close()
	time.Sleep(1 * time.Second)
	// close the file
	if err := f.Close(); err != nil {
		t.Fatalf("unable to close the file: %s", f.Name())
	}

	// and remove the file
	if err := os.Remove(f.Name()); err != nil {
		t.Fatalf("unable to remove the file from %s", f.Name())
	}

	time.Sleep(1 * time.Second)
}
```
> `store.go`
```go
package main

import (
	"bytes"

	bolt "github.com/etcd-io/bbolt"
)
//Panic，如果您不想因严重的INITIALIZE-ONLY-ERRORS 异常而将其更改

// Panic panics, change it if you don't want to panic on critical INITIALIZE-ONLY-ERRORS
var Panic = func(v interface{}) {
	panic(v)
}

// Store是网址的存储接口
//注意：没有Del功能

// Store is the store interface for urls.
// Note: no Del functionality.
type Store interface {
	Set(key string, value string) error // 如果出了问题返回错误 | error if something went wrong
	Get(key string) string              // 如果找不到，则为空值 | empty value if not found
	Len() int                           // 应该返回所有记录/表/桶的数量 | should return the number of all the records/tables/buckets
	Close()                             // 释放存储或忽略 | release the store or ignore
}

var tableURLs = []byte("urls")

//Store的数据库表示形式。
//只有一个表/存储桶包含网址，因此它不是完整的数据库，
//它仅适用于单个存储桶，因为我们需要这些。

// DB representation of a Store.
// Only one table/bucket which contains the urls, so it's not a fully Database,
// it works only with single bucket because that all we need.
type DB struct {
	db *bolt.DB
}

var _ Store = &DB{}
// openDatabase打开一个新的数据库连接并返回其实例

// openDatabase open a new database connection
// and returns its instance.
func openDatabase(stumb string) *bolt.DB {
	//打开当前工作目录下的data（base）文件，
	// 如果不存在该文件将被创建

	// Open the data(base) file in the current working directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(stumb, 0600, nil)
	if err != nil {
		Panic(err)
	}
	//在此处创建存储桶

	// create the buckets here
	tables := [...][]byte{
		tableURLs,
	}

	db.Update(func(tx *bolt.Tx) (err error) {
		for _, table := range tables {
			_, err = tx.CreateBucketIfNotExists(table)
			if err != nil {
				Panic(err)
			}
		}

		return
	})

	return db
}

// NewDB返回一个新的数据库实例，其连接已打开
// DB实现Store

// NewDB returns a new DB instance, its connection is opened.
// DB implements the Store.
func NewDB(stumb string) *DB {
	return &DB{
		db: openDatabase(stumb),
	}
}
// Set设置一个缩短的网址及其键
//注意：调用方负责生成密钥

// Set sets a shorten url and its key
// Note: Caller is responsible to generate a key.
func (d *DB) Set(key string, value string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(tableURLs)
		//生成url的ID注意：我们可以使用它代替随机的字符串键
		// 但是我们想模拟一个实际的url缩短器，因此我们跳过它
		// id, _ := b.NextSequence()

		// Generate ID for the url
		// Note: we could use that instead of a random string key
		// but we want to simulate a real-world url shortener
		// so we skip that.
		// id, _ := b.NextSequence()
		if err != nil {
			return err
		}

		k := []byte(key)
		valueB := []byte(value)
		c := b.Cursor()

		found := false
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Equal(valueB, v) {
				found = true
				break
			}
		}
		//如果值已经存在，请不要重新输入

		// if value already exists don't re-put it.
		if found {
			return nil
		}

		return b.Put(k, []byte(value))
	})
}
//Clear将清除表URL的所有数据库条目

// Clear clears all the database entries for the table urls.
func (d *DB) Clear() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(tableURLs)
	})
}
// Get通过其键返回一个URL
//
//如果找不到，则返回一个空字符串

// Get returns a url by its key.
//
// Returns an empty string if not found.
func (d *DB) Get(key string) (value string) {
	keyB := []byte(key)
	d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(tableURLs)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Equal(keyB, k) {
				value = string(v)
				break
			}
		}

		return nil
	})

	return
}

// GetByValue返回特定（原始）URL值的所有键

// GetByValue returns all keys for a specific (original) url value.
func (d *DB) GetByValue(value string) (keys []string) {
	valueB := []byte(value)
	d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(tableURLs)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		//首先为存储区的表"urls"
		// first for the bucket's table "urls"
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Equal(valueB, v) {
				keys = append(keys, string(k))
			}
		}

		return nil
	})

	return
}
// Len返回所有“短”网址的长度
// Len returns all the "shorted" urls length
func (d *DB) Len() (num int) {
	d.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(tableURLs)
		if b == nil {
			return nil
		}

		b.ForEach(func([]byte, []byte) error {
			num++
			return nil
		})
		return nil
	})
	return
}
//关闭将关闭数据库）连接

// Close shutdowns the data(base) connection.
func (d *DB) Close() {
	if err := d.db.Close(); err != nil {
		Panic(err)
	}
}
```