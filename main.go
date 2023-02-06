package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

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
	log.Println("Connect to", r.RemoteAddr)
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

func getMyIPV6() string {
	s, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range s {
		i := regexp.MustCompile(`(\w+:){7}\w+`).FindString(a.String())
		if strings.Count(i, ":") == 7 {
			return i
		}
	}
	return ""
}

func main() {
	host := "[" + getMyIPV6() + "]:8002"
	log.Printf("base url http://%s\n", host)

	mux := http.NewServeMux()

	files := http.FileServer(http.Dir("./"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.HandleFunc("/", getInitPage)
	mux.HandleFunc("/ws", updateWebsocket)

	http.ListenAndServe(host, mux)
}
