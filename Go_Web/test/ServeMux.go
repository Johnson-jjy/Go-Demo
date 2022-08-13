package main

import (
	"fmt"
	"net/http"
	"os"
)

type WelcomeHandlerStruct struct {

}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
	fmt.Println("HW")
	os.PathError{}
	os.ErrClosed
}

func (*WelcomeHandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome")
	fmt.Println("HWW")
}

func main () {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloHandler)
	mux.Handle("/welcome", &WelcomeHandlerStruct{})
	http.ListenAndServe(":8181", mux)
}
