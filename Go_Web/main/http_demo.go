package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func sayHelloName(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm() // 解析参数,默认是不会解析的
	fmt.Println(r.Form) // 这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Johnson!") // 这个写入到w的回输出到客户端
}

func main()  {
	http.HandleFunc("/", sayHelloName) // 设置访问的路由(即默认路由)
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口; nil->匹配默认路由(DefaultServeMux)
	// DefaultServeMux会调用ServeHTTP,其内部调用sayHelloName本身
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
