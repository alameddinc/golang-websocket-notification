package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000

	MaxBufferSize = 4096
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  MaxBufferSize,
	WriteBufferSize: MaxBufferSize,
}

// Client represents the websocket client at the server
type Client struct {
	// The actual websocket connection.
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
}

func newClient(conn *websocket.Conn, server *WsServer, name string) *Client {
	return &Client{
		ID:       uuid.New(),
		Name:     name,
		conn:     conn,
		wsServer: server,
		send:     make(chan []byte, 256),
	}
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(ws *WsServer, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	if !ok {
		log.Println("Url Param 'name' is missing")
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, ws, name[0])
	go client.writePump()
	go client.readPump()
	ws.register <- client
}

func SendMessageAllUsers(message []byte) {
	WS.broadcastToClient(message)
}

func SendMessageUser(userUUID string, message []byte) {
	fmt.Printf(".")
	if client := WS.findClient(userUUID); client != nil {
		fmt.Printf("Yollandı \n")
		client.send <- message
	}
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

	}

}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) disconnect() {
	c.wsServer.unregister <- c
	close(c.send)
	c.conn.Close()
}

func (c *Client) getName() string {
	return c.Name
}
