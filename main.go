package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
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
	    <script>
		let sockAddr = 'ws://' + window.location.host + '/ws';
		let socket = new WebSocket(sockAddr);
		console.log("Attempting Connection...");

		socket.onopen = () => {
		    console.log("Successfully Connected");
		    socket.send("Hi From the Client!")
		};
		
		socket.onclose = event => {
		    console.log("Socket Closed Connection: ", event);
		    socket.send("Client Closed!")
		};

		socket.onerror = error => {
		    console.log("Socket Error: ", error);
		};

	    </script>
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
	http.HandleFunc("/ws", webSocketHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
