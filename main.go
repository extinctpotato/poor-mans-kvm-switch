package main

import (
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

		fmt.Println(string(p))
		s.Write(p)

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
	    <canvas width="600" height="600">
	    </canvas>
	    <script>
		let sockAddr = 'ws://' + window.location.host + '/ws';
		let socket = new WebSocket(sockAddr);

		// Get canvas-related DOM variables.
		let canvas = document.querySelector('canvas');
		let ctx = canvas.getContext('2d');

		// Draw black background.
		ctx.fillStyle = "black";
		ctx.fillRect(0, 0, canvas.width, canvas.height);
		ctx.fill();

		// Allow grabbing cursor and releasing it.
		canvas.onclick = function() {
		  canvas.requestPointerLock();
		}

		function sendMovement(e) {
		  let a = Math.min(Math.abs(e.movementX) + 1, 127);
		  let b = Math.min(Math.abs(e.movementY) + 1, 127);

		  let c = [e.movementX > 0, e.movementY > 0].reduce((res, x) => res << 1 | x);
		  socket.send(String.fromCharCode(65, a, b, c));
		}

		function sendMouseClick(e) {
		  let msgType = e.type === 'mouseup' ? 68 : 67;
		  socket.send(String.fromCharCode(msgType, e.button, 0, 0));
		}

		document.addEventListener("pointerlockchange", () => {
		  if (document.pointerLockElement === canvas) {
		    document.addEventListener('mousemove', sendMovement);
		    document.addEventListener('mouseup', sendMouseClick);
		    document.addEventListener('mousedown', sendMouseClick);
		  } else {
		    document.removeEventListener('mousemove', sendMovement);
		    document.removeEventListener('mouseup', sendMouseClick);
		    document.removeEventListener('mousedown', sendMouseClick);
		  }
		});

		console.log("Attempting Connection...");

		socket.onopen = () => {
		    console.log("Successfully Connected");
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
