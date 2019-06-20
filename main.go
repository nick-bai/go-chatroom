package main

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"html/template"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter,r *http.Request)  {

	tem,err := template.ParseFiles("view/index.html")
	if err != nil{
		fmt.Println("读取文件失败,err",err)
		return
	}

	tem.Execute(w, "")
}

func main()  {

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "您的答案是： " + msg
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/chat", indexHandler)

	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
