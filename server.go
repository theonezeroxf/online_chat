package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/chat_m.html") // 加载模板文件
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil) // 渲染模板并返回给客户端
}
func main() {
	router := mux.NewRouter()
	go h.run()
	router.HandleFunc("/ws", myws)
	router.HandleFunc("/chat", chatHandler)
	fmt.Println("Start Server at 0.0.0.0:8080")
	http.ListenAndServe(":8080", router)
	// fmt.Println("Start Server at 0.0.0.0:8080")
}
