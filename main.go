package main

import (
	"fmt"
	"github.com/googollee/go-engine.io"
	"github.com/googollee/go-engine.io/transport"
	"github.com/googollee/go-engine.io/transport/polling"
	"github.com/googollee/go-engine.io/transport/websocket"
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

	pt := polling.Default

	wt := websocket.Default
	wt.CheckOrigin = func(req *http.Request) bool {
		return true
	}

	server, err := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			pt,
			wt,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
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
