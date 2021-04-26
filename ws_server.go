package main

import "fmt"

type WsServer struct {
	Clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewWSServer() *WsServer {
	return &WsServer{
		Clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (w *WsServer) Run() {
	for {
		select {
		case client := <-w.register:
			fmt.Println("Bağlantı Kuruldu")
			w.RegisterFunction(client)
		case client := <-w.unregister:
			fmt.Println("Bağlantı Koptu")
			w.UnRegisterFunction(client)
		case message := <-w.broadcast:
			fmt.Println("Mesaj alındı")
			w.broadcastToClient(message)
		}
	}
}

func (w *WsServer) RegisterFunction(client *Client) {
	w.Clients[client] = true
}

func (w *WsServer) UnRegisterFunction(client *Client) {
	if _, ok := w.Clients[client]; ok {
		w.Clients[client] = false
		delete(w.Clients, client)
	}
}

func (server *WsServer) findClient(c string) *Client {
	var foundClient *Client
	for client, status := range server.Clients {
		if c == client.Name {
			if status {
				foundClient = client
			}
			break
		}
	}
	return foundClient
}

func (w *WsServer) broadcastToClient(message []byte) {
	for client := range w.Clients {
		client.send <- message
	}
}
