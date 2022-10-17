package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/learningspoons-go/groupchat/api"
	"github.com/learningspoons-go/groupchat/db"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	if err := db.Connect("postgres://postgres:chatdbpasswd123@chatdb.learningspoons.danslee.com:5432/postgres"); err != nil {
		log.Fatalln(err)
	}

	if err := http.ListenAndServe(":8080", api.Handler()); err != nil {
		log.Fatalln(err)
	}
}
