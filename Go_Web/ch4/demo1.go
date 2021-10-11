package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析url传递的参数,对于POST则解析响应包的主体 ( request body )
	// 注: 若没有调用ParseForm方法, 下面无法获取表单的数据
	fmt.Println(r.Form) // 这些信息是输出到服务器端的打印信息
	fmt.Println("path:", r.URL.Path)
	fmt.Println("scheme:", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Johnson!") // 这个是写入到w的是输出给客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("D:\\go_WorkSpace\\src\\Go-Demo\\Go_Web\\ch4\\login.gtml")
		t.Execute(w, nil)
	} else {
		// 请求的是登陆数据, 那么执行登录的逻辑判断
		r.ParseForm() // 默认情况下,Handler不会自动解析form,显式调用此行语句才能对表单数据进行操作
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
	// 可对form数据进行的操作如下:
	v := url.Values{}
	v.Set("name", "zyx")
	v.Add("friend", "hb")
	v.Add("friend", "hl")
	v.Add("friend", "shl")
	fmt.Println(v.Encode())
	fmt.Println(v.Get("name"))
	fmt.Println(v.Get("friend"))
	fmt.Println(v["friend"])
}

func main()  {
	http.HandleFunc("/", sayHelloName) // 设置访问的路由
	http.HandleFunc("/login", login) // 设置访问的路由
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

