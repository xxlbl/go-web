package main

// $ curl http://localhost:9999/
// URL.Path = "/"
// $ curl http://localhost:9999/hello
// Header["Accept"] = ["*/*"]
// Header["User-Agent"] = ["curl/7.54.0"]

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//设置了2个路由，/和/hello，分别绑定 indexHandler 和 helloHandler
	//根据不同的HTTP请求会调用不同的处理函数。
	//访问/，响应是URL.Path = /，
	//而/hello的响应则是请求头(header)中的键值对信息。
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	//启动 Web 服务，
	//第一个参数是地址，:9999表示在 9999 端口监听。
	//而第二个参数则代表处理所有的HTTP请求的实例，nil 代表使用*标准库*中的实例处理。
	log.Fatal(http.ListenAndServe(":9999", nil))
}

//第二个参数Handler是一个接口，需要实现方法 ServeHTTP ，
//只要传入任何实现了 ServerHTTP 接口的实例，所有的HTTP请求，就都交给了该实例处理了
//package http
//
//type Handler interface {
//    ServeHTTP(w ResponseWriter, r *Request)
//}
//
//func ListenAndServe(address string, h Handler) error

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
