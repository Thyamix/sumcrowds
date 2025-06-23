package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thyamix/festival-counter/internal/models"
)

type Client struct {
	Server       *Server
	Conn         *websocket.Conn
	FestivalCode string
	Send         chan []byte
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Code    string          `json:"code"`
	Content json.RawMessage `json:"content"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Client) readPump() {
	defer func() {
		c.Server.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Read error:", err)
			}
			break
		}

		handleMessage(c, msg)
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func HandleCounter(server *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	festivalCode := r.PathValue("festivalCode")
	fmt.Printf("Starting WS connection on %v \n", festivalCode)
	if err != nil {
		log.Println("Upgader error: ", err)
		return
	}

	client := &Client{Server: server, Conn: conn, Send: make(chan []byte), FestivalCode: festivalCode}
	client.Server.Register <- client

	go client.readPump()
	go client.writePump()
}

func handleMessage(c *Client, message []byte) error {
	var incomingMessage IncomingMessage
	if err := json.Unmarshal(message, &incomingMessage); err != nil {
		return fmt.Errorf("invalid message from handle message: %v", err)
	}

	switch incomingMessage.Type {
	case "ping":
		pong(c)
	case "getTotal":
		sendTotal(c)
	}
	return nil
}

func sendTotal(c *Client) error {
	err := BroadcastTotal(c.FestivalCode)
	return err
}

func pong(c *Client) {
	pingJson, err := json.Marshal(models.Response{Type: "pong"})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("pong")
	c.Send <- pingJson
}
