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
