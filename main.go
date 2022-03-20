package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tarm/serial"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//go:embed index.js
var indexJs string

func reader(conn *websocket.Conn) {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200, Size: 8}

	s, err := serial.OpenPort(c)

	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(p)
		s.Write(p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func indexJsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, indexJs)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	  <head>
	    <meta charset="UTF-8" />
	    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
	    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
	    <title>poor mans kvm test</title>
	  </head>
	  <body>
	    <canvas width="600" height="600"></canvas>
	    <script src="index.js"></script>
	  </body>
	</html>
	`

	fmt.Fprintf(w, html)
}

func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement this properly.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	// TODO: handle the error properly.
	if err != nil {
		log.Println(err)
	}

	err = ws.WriteMessage(1, []byte("test"))
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

func main() {
	fmt.Println("vim-go")

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/index.js", indexJsHandler)
	http.HandleFunc("/ws", webSocketHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
