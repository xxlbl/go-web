package main

// $ curl http://localhost:9999/
// URL.Path = "/"
// $ curl http://localhost:9999/hello
// Header["Accept"] = ["*/*"]
// Header["User-Agent"] = ["curl/7.54.0"]
// curl http://localhost:9999/world
// 404 NOT FOUND: /world

import (
	"fmt"
	"log"
	"net/http"
)

// Engine is the uni handler for all requests
type Engine struct{}

//定义了一个空的结构体Engine，实现了方法ServeHTTP
//第二个参数是 Request ，该对象包含了*该HTTP请求的所有的信息*，比如请求地址、Header和Body等信息；
//第一个参数是 ResponseWriter ，利用 ResponseWriter 可以构造针对该请求的响应。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	//web框架实现第一步：将所有的HTTP请求转向了我们自己的处理逻辑
	log.Fatal(http.ListenAndServe(":9999", engine))
}

//在实现Engine之前，调用 http.HandleFunc 实现了路由和Handler的映射，
//也就是只能针对具体的路由写处理逻辑。比如/hello。
//但是在实现Engine之后，我们拦截了所有的HTTP请求，拥有了统一的控制入口。
