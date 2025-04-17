package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 处理/hello路径的请求
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "你好，Gopher！")
	})

	// 启动服务器在8080端口
	fmt.Println("服务器启动在 :8080...")
	http.ListenAndServe(":8080", nil)
}
