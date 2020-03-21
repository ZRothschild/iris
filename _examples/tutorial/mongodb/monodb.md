# 使用官方的MongoDB Go驱动程序和Iris构建RESTful API

文章即将发布，关注并继续关注

- <https://medium.com/@kataras>
- <https://dev.to/kataras>

查看 [功能齐全的例子](main.go).

## 目录结构
> 主目录`mongodb`
```html
    —— api
        —— store
            —— movie.go
    —— env
        —— env.go
    —— httputil
        —— error.go
    —— store
        —— movie.go
    —— .env
    —— main.go
```

## 代码示例
> `api/store/movie.go`
```golang
package storeapi

import (
	"github.com/kataras/iris/v12/_examples/tutorial/mongodb/httputil"
	"github.com/kataras/iris/v12/_examples/tutorial/mongodb/store"

	"github.com/kataras/iris/v12"
)

type MovieHandler struct {
	service store.MovieService
}

func NewMovieHandler(service store.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (h *MovieHandler) GetAll(ctx iris.Context) {
	movies, err := h.service.GetAll(nil)
	if err != nil {
		httputil.InternalServerErrorJSON(ctx, err, "Server was unable to retrieve all movies")
		return
	}

	if movies == nil {
		//如果movies为空，则将返回"null"，就可以使用此“技巧”，将null 转换成"[]"空数组json返回
		
		// will return "null" if empty, with this "trick" we return "[]" json.
		movies = make([]store.Movie, 0)
	}

	ctx.JSON(movies)
}

func (h *MovieHandler) Get(ctx iris.Context) {
	id := ctx.Params().Get("id")

	m, err := h.service.GetByID(nil, id)
	if err != nil {
		if err == store.ErrNotFound {
			ctx.NotFound()
		} else {
			httputil.InternalServerErrorJSON(ctx, err, "Server was unable to retrieve movie [%s]", id)
		}
		return
	}

	ctx.JSON(m)
}

func (h *MovieHandler) Add(ctx iris.Context) {
	m := new(store.Movie)

	err := ctx.ReadJSON(m)
	if err != nil {
		httputil.FailJSON(ctx, iris.StatusBadRequest, err, "Malformed request payload")
		return
	}

	err = h.service.Create(nil, m)
	if err != nil {
		httputil.InternalServerErrorJSON(ctx, err, "Server was unable to create a movie")
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(m)
}

func (h *MovieHandler) Update(ctx iris.Context) {
	id := ctx.Params().Get("id")

	var m store.Movie
	err := ctx.ReadJSON(&m)
	if err != nil {
		httputil.FailJSON(ctx, iris.StatusBadRequest, err, "Malformed request payload")
		return
	}

	err = h.service.Update(nil, id, m)
	if err != nil {
		if err == store.ErrNotFound {
			ctx.NotFound()
			return
		}
		httputil.InternalServerErrorJSON(ctx, err, "Server was unable to update movie [%s]", id)
		return
	}
}

func (h *MovieHandler) Delete(ctx iris.Context) {
	id := ctx.Params().Get("id")

	err := h.service.Delete(nil, id)
	if err != nil {
		if err == store.ErrNotFound {
			ctx.NotFound()
			return
		}
		httputil.InternalServerErrorJSON(ctx, err, "Server was unable to delete movie [%s]", id)
		return
	}
}
```
> `env/env.go`
```golang
package env

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var (
	// Port是PORT环境变量，如果缺少，则为8080。
	//用于打开Web服务器的tcp侦听器。

	// Port is the PORT environment variable or 8080 if missing.
	// Used to open the tcp listener for our web server.
	Port string
	// DSN是DSN环境变量，如果缺少，则为mongodb：// localhost：27017。
	//用于连接到mongodb。

	// DSN is the DSN environment variable or mongodb://localhost:27017 if missing.
	// Used to connect to the mongodb.
	DSN string
)

func parse() {
	Port = getDefault("PORT", "8080")
	DSN = getDefault("DSN", "mongodb://localhost:27017")
}

//Load加载整个应用程序中正在使用的环境变量
//从文件加载，即.env或dev.env

// Load loads environment variables that are being used across the whole app.
// Loading from file(s), i.e .env or dev.env
//
// Example of a 'dev.env':
// PORT=8080
// DSN=mongodb://localhost:27017

//在加载之后，调用者可以通过os.Getenv获取环境变量

// After `Load` the callers can get an environment variable via `os.Getenv`.
func Load(envFileName string) {
	if args := os.Args; len(args) > 1 && args[1] == "help" {
		fmt.Fprintln(os.Stderr, "https://github.com/kataras/iris/blob/master/_examples/tutorials/mongodb/README.md")
		os.Exit(-1)
	}

	log.Printf("Loading environment variables from file: %s\n", envFileName)
	//如果有多个文件名以逗号分隔，则从所有文件名中加载，所以每一个env文件都可以放特定的东西

	// If more than one filename passed with comma separated then load from all
	// of these, a env file can be a partial too.
	envFiles := strings.Split(envFileName, ",")
	for i := range envFiles {
		if filepath.Ext(envFiles[i]) == "" {
			envFiles[i] += ".env"
		}
	}

	if err := godotenv.Load(envFiles...); err != nil {
		panic(fmt.Sprintf("error loading environment variables from [%s]: %v", envFileName, err))
	}

	envMap, _ := godotenv.Read(envFiles...)
	for k, v := range envMap {
		log.Printf("◽ %s=%s\n", k, v)
	}

	parse()
}

func getDefault(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		os.Setenv(key, def)
		value = def
	}

	return value
}
```
> `httputil/error.go`
```golang
package httputil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

var validStackFuncs = []func(string) bool{
	func(file string) bool {
		return strings.Contains(file, "/mongodb/api/")
	},
}
// RuntimeCallerStack返回应用程序的`file:line`(错误所在文件的行数) 堆栈跟踪，以提供有关错误原因的更多信息

// RuntimeCallerStack returns the app's `file:line` stacktrace
// to give more information about an error cause.
func RuntimeCallerStack() (s string) {
	var pcs [10]uintptr
	n := runtime.Callers(1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		for _, fn := range validStackFuncs {
			if fn(frame.File) {
				s += fmt.Sprintf("\n\t\t\t%s:%d", frame.File, frame.Line)
			}
		}

		if !more {
			break
		}
	}

	return s
}
// HTTPError描述HTTP错误

// HTTPError describes an HTTP error.
type HTTPError struct {
	error
	Stack       string    `json:"-"` // 整个堆栈跟踪 | the whole stacktrace.
	CallerStack string    `json:"-"` // 调用者，文件：lineNumber | the caller, file:lineNumber
	When        time.Time `json:"-"` // 错误发生的时间 | the time that the error occurred.
	// ErrorCode int：可能是已知错误代码的集合

	// ErrorCode int: maybe a collection of known error codes.
	StatusCode int `json:"statusCode"`
	//也可以命名为“原因”，它也是错误的信息

	// could be named as "reason" as well
	//  it's the message of the error.
	Description string `json:"description"`
}

func newError(statusCode int, err error, format string, args ...interface{}) HTTPError {
	if format == "" {
		format = http.StatusText(statusCode)
	}

	desc := fmt.Sprintf(format, args...)
	if err == nil {
		err = errors.New(desc)
	}

	return HTTPError{
		err,
		string(debug.Stack()),
		RuntimeCallerStack(),
		time.Now(),
		statusCode,
		desc,
	}
}

func (err HTTPError) writeHeaders(ctx iris.Context) {
	ctx.StatusCode(err.StatusCode)
	ctx.Header("X-Content-Type-Options", "nosniff")
}
// LogFailure将失败输出到"logger"

// LogFailure will print out the failure to the "logger".
func LogFailure(logger io.Writer, ctx iris.Context, err HTTPError) {
	timeFmt := err.When.Format("2006/01/02 15:04:05")
	firstLine := fmt.Sprintf("%s %s: %s", timeFmt, http.StatusText(err.StatusCode), err.Error())
	whitespace := strings.Repeat(" ", len(timeFmt)+1)
	fmt.Fprintf(logger, "%s\n%sIP: %s\n%sURL: %s\n%sSource: %s\n",
		firstLine, whitespace, ctx.RemoteAddr(), whitespace, ctx.FullRequestURI(), whitespace, err.CallerStack)
}

//失败将发送状态码，写出错误原因
//并返回HTTPError供进一步使用，即记录，请参见`InternalServerError`。

// Fail will send the status code, write the error's reason
// and return the HTTPError for further use, i.e logging, see `InternalServerError`.
func Fail(ctx iris.Context, statusCode int, err error, format string, args ...interface{}) HTTPError {
	httpErr := newError(statusCode, err, format, args...)
	httpErr.writeHeaders(ctx)

	ctx.WriteString(httpErr.Description)
	return httpErr
}

// FailJSON将错误数据作为JSON发送给客户端。
//对于API很有用。

// FailJSON will send to the client the error data as JSON.
// Useful for APIs.
func FailJSON(ctx iris.Context, statusCode int, err error, format string, args ...interface{}) HTTPError {
	httpErr := newError(statusCode, err, format, args...)
	httpErr.writeHeaders(ctx)

	ctx.JSON(httpErr)

	return httpErr
}
// InternalServerError记录到服务器的终端，并将500 Internal Server Error分发给客户端
// 内部服务器错误至关重要，因此我们将其记录到os.Stderr中

// InternalServerError logs to the server's terminal
// and dispatches to the client the 500 Internal Server Error.
// Internal Server errors are critical, so we log them to the `os.Stderr`.
func InternalServerError(ctx iris.Context, err error, format string, args ...interface{}) {
	LogFailure(os.Stderr, ctx, Fail(ctx, iris.StatusInternalServerError, err, format, args...))
}
// InternalServerErrorJSON的行为与`InternalServerError`完全相同，但是它将数据作为JSON发送
//对于API很有用

// InternalServerErrorJSON acts exactly like `InternalServerError` but instead it sends the data as JSON.
// Useful for APIs.
func InternalServerErrorJSON(ctx iris.Context, err error, format string, args ...interface{}) {
	LogFailure(os.Stderr, ctx, FailJSON(ctx, iris.StatusInternalServerError, err, format, args...))
}
// UnauthorizedJSON发送StatusUnauthorized（401）HTTPError值的JSON格式

// UnauthorizedJSON sends JSON format of StatusUnauthorized(401) HTTPError value.
func UnauthorizedJSON(ctx iris.Context, err error, format string, args ...interface{}) HTTPError {
	return FailJSON(ctx, iris.StatusUnauthorized, err, format, args...)
}
```
> `store/movie.go`
```golang
package store

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// up to you:
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type Movie struct {
	/*您需要 bson:"_id" 才能以ID填充检索 */
	ID          primitive.ObjectID `json:"_id" bson:"_id"` /* you need the bson:"_id" to be able to retrieve with ID filled */
	Name        string             `json:"name"`
	Cover       string             `json:"cover"`
	Description string             `json:"description"`
}

type MovieService interface {
	GetAll(ctx context.Context) ([]Movie, error)
	GetByID(ctx context.Context, id string) (Movie, error)
	Create(ctx context.Context, m *Movie) error
	Update(ctx context.Context, id string, m Movie) error
	Delete(ctx context.Context, id string) error
}

type movieService struct {
	C *mongo.Collection
}

var _ MovieService = (*movieService)(nil)

func NewMovieService(collection *mongo.Collection) MovieService {
	// up to you:
	// indexOpts := new(options.IndexOptions)
	// indexOpts.SetName("movieIndex").
	// 	SetUnique(true).
	// 	SetBackground(true).
	// 	SetSparse(true)

	// collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
	// 	Keys:    []string{"_id", "name"},
	// 	Options: indexOpts,
	// })

	return &movieService{C: collection}
}

func (s *movieService) GetAll(ctx context.Context) ([]Movie, error) {
	// 注意：
	// mongodb的go-driver文档中，您可以将`nil`传递给"find all"，但这会导致NilDocument错误，可能是错误或documentation错误
	// 您必须传递`bson.D {}`。

	// Note:
	// The mongodb's go-driver's docs says that you can pass `nil` to "find all" but this gives NilDocument error,
	// probably it's a bug or a documentation's mistake, you have to pass `bson.D{}` instead.
	cur, err := s.C.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []Movie

	for cur.Next(ctx) {
		if err = cur.Err(); err != nil {
			return nil, err
		}

		//	elem := bson.D{}
		var elem Movie
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		// results = append(results, Movie{ID: elem[0].Value.(primitive.ObjectID)})

		results = append(results, elem)
	}

	return results, nil
}

func matchID(id string) (bson.D, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	return filter, nil
}

var ErrNotFound = errors.New("not found")

func (s *movieService) GetByID(ctx context.Context, id string) (Movie, error) {
	var movie Movie
	filter, err := matchID(id)
	if err != nil {
		return movie, err
	}

	err = s.C.FindOne(ctx, filter).Decode(&movie)
	if err == mongo.ErrNoDocuments {
		return movie, ErrNotFound
	}
	return movie, err
}

func (s *movieService) Create(ctx context.Context, m *Movie) error {
	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	_, err := s.C.InsertOne(ctx, m)
	if err != nil {
		return err
	}
	//如果Movie.ID字段上有`bson:"_id`，以下内容将不起作用，
	//没有`bson:"_id`，我们就需要手动生成了一个新ID（如上所示）

	// The following doesn't work if you have the `bson:"_id` on Movie.ID field,
	// therefore we manually generate a new ID (look above).
	// res, err := ...InsertOne
	// objectID := res.InsertedID.(primitive.ObjectID)
	// m.ID = objectID
	return nil
}

func (s *movieService) Update(ctx context.Context, id string, m Movie) error {
	filter, err := matchID(id)
	if err != nil {
		return err
	}

	// update := bson.D{
	// 	{Key: "$set", Value: m},
	// }
	// ^这将覆盖所有字段，您可以执行此操作，具体取决于您的设计。 但是让我们检查每个字段：
	
	// ^ this will override all fields, you can do that, depending on your design. but let's check each field:
	elem := bson.D{}

	if m.Name != "" {
		elem = append(elem, bson.E{Key: "name", Value: m.Name})
	}

	if m.Description != "" {
		elem = append(elem, bson.E{Key: "description", Value: m.Description})
	}

	if m.Cover != "" {
		elem = append(elem, bson.E{Key: "cover", Value: m.Cover})
	}

	update := bson.D{
		{Key: "$set", Value: elem},
	}

	_, err = s.C.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *movieService) Delete(ctx context.Context, id string) error {
	filter, err := matchID(id)
	if err != nil {
		return err
	}
	_, err = s.C.DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return err
	}

	return nil
}
```
> `.env`
```sh
# .env file contents
PORT=8080
DSN=localhost:27017
```
> `main.go`
```golang
package main

// go get -u go.mongodb.org/mongo-driver
// go get -u github.com/joho/godotenv

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	// APIs
	storeapi "github.com/kataras/iris/v12/_examples/tutorial/mongodb/api/store"

	//
	"github.com/kataras/iris/v12/_examples/tutorial/mongodb/env"
	"github.com/kataras/iris/v12/_examples/tutorial/mongodb/store"

	"github.com/kataras/iris/v12"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const version = "0.0.1"

func init() {
	envFileName := ".env"

	flagset := flag.CommandLine
	flagset.StringVar(&envFileName, "env", envFileName, "the env file which web app will use to extract its environment variables")
	flag.CommandLine.Parse(os.Args[1:])

	env.Load(envFileName)
}

func main() {
	clientOptions := options.Client().SetHosts([]string{env.DSN})
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("store")

	var (
		// 集合 | Collections.
		moviesCollection = db.Collection("movies")

		// Services.
		movieService = store.NewMovieService(moviesCollection)
	)

	app := iris.New()
	app.Use(func(ctx iris.Context) {
		ctx.Header("Server", "Iris MongoDB/"+version)
		ctx.Next()
	})

	storeAPI := app.Party("/api/store")
	{
		movieHandler := storeapi.NewMovieHandler(movieService)
		storeAPI.Get("/movies", movieHandler.GetAll)
		storeAPI.Post("/movies", movieHandler.Add)
		storeAPI.Get("/movies/{id}", movieHandler.Get)
		storeAPI.Put("/movies/{id}", movieHandler.Update)
		storeAPI.Delete("/movies/{id}", movieHandler.Delete)
	}

	// GET: http://localhost:8080/api/store/movies
	// POST: http://localhost:8080/api/store/movies
	// GET: http://localhost:8080/api/store/movies/{id}
	// PUT: http://localhost:8080/api/store/movies/{id}
	// DELETE: http://localhost:8080/api/store/movies/{id}
	app.Run(iris.Addr(fmt.Sprintf(":%s", env.Port)), iris.WithOptimizations)
}
```

## 图片
### 添加movie
![https://studyiris.com](0_create_movie.png)

### 更新movie

![https://studyiris.com](1_update_movie.png)

### 获取所有movie

![https://studyiris.com](2_get_all_movies.png)

### 通过Id获取特定movie

![https://studyiris.com](3_get_movie.png)

### 根据Id删除特定movie

![https://studyiris.com](4_delete_movie.png)

### mongodb 安装 【我的是Ubuntu18】
1. 下载地址[mongodb](https://www.mongodb.com/download-center/community) 下载`$ sudo wget https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1804-4.2.3.tgz`
2. 解压安装包 `sudo tar -zvxf mongodb-linux-x86_64-ubuntu1804-4.2.3.tgz`
3. 编辑 `/etc/environment` 加入`/usr/local/soft/mongodb-linux-x86_64-ubuntu1804-4.2.3/bin`【看你自己的地址这是我的哦大佬】
4. 创建 `/etc/monogdb.cof`  `/var/log/mongodb/mongodb.log`   `/data/mongodb/db` 需要特别注意权限
5. `mongod --config /etc/monogdb.cof --fork`  **--fork**生成进程
```editorconfig
# 日志文件位置
logpath=/var/log/mongodb/mongodb.log

#以追加方式写入日志
logappend=true

# 默认27017
#port= 27017

# 是否以守护进程方式运行
#fork=true

# 数据库文件位置
dbpath=/data/mongodb/db

# 启用定期记录CPU利用率和 I/O 等待
#cpu=true

# 是否以安全认证方式运行，默认是不认证的非安全方式
#noauth=true
#auth=true
#
bind_ip=0.0.0.0
```
