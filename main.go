package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func getInitPage(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("InitPage.html")
	if err != nil {
		log.Panicf("Failed to open \"InitPage.html\": %v\n", err)
	}
	buf_f := bufio.NewReader(f)
	_, err = io.Copy(w, buf_f)
	if err != nil {
		log.Panicf("Failed to send \"InitPage.html\": %v\n", err)
	}
}

func updateWebsocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Host)
	fmt.Println(r.Host)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	go handleWebsocketConn(ws)
}

func handleWebsocketConn(ws *websocket.Conn) {
	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := ws.WriteMessage(messageType, msg); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	mux := http.NewServeMux()

	files := http.FileServer(http.Dir("./"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.HandleFunc("/", getInitPage)
	mux.HandleFunc("/ws", updateWebsocket)

	http.ListenAndServe("192.168.108.3:8002", mux)
}
