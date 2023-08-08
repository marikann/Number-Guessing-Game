package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// WebSocket bağlantısı struct yapısı
type WebSocketConnection struct {
	conn *websocket.Conn
	id   string
}

var connections = make(map[*WebSocketConnection]bool)

// WebSocket bağlantısı için upgrader, tüm istekleri kabul etsin diye checkorgin i true dönüyorum
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
