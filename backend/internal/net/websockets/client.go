package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thyamix/festival-counter/internal/database"
)

type Client struct {
	Server       *Server
	Conn         *websocket.Conn
	FestivalCode string
	Send         chan []byte
}

type ValueChange struct {
	Amount int `json:"amount"`
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
	case "inc":
		inc(c, incomingMessage.Content)
	case "dec":
		dec(c, incomingMessage.Content)
	case "getTotal":
		sendTotal(c, false)
	}
	return nil
}

func sendTotal(c *Client, broadcast bool) error {
	total, maxGauge, err := database.GetTotal(c.FestivalCode)
	if err != nil {
		log.Print(err)
	}
	totalJson, err := json.Marshal(map[string]int{
		"total": total,
		"jauge": maxGauge,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Sending total:max: %v:%v \n", total, maxGauge)

	if broadcast {
		c.Server.Messages <- Message{Message: totalJson, FestivalCode: c.FestivalCode}
	} else {
		c.Send <- totalJson
	}
	return nil
}

func inc(c *Client, data json.RawMessage) error {
	var valueChange ValueChange

	err := json.Unmarshal(data, &valueChange)

	if err != nil {
		return err
	}

	amount := valueChange.Amount

	if amount <= 0 || amount > 100 {
		return fmt.Errorf("can't increment by %v as it is not between 0 - 100", amount)
	}

	total, _, err := database.GetTotal(c.FestivalCode)
	if err != nil {
		fmt.Println(err)
	}

	err = database.AddValue(amount, c.FestivalCode)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Value change on: ", c.FestivalCode)
	fmt.Println("+", amount)
	fmt.Println("New total of", total+amount)

	sendTotal(c, true)
	return nil
}

func dec(c *Client, data json.RawMessage) error {
	var valueChange ValueChange

	err := json.Unmarshal(data, &valueChange)
	if err != nil {
		return err
	}

	amount := valueChange.Amount

	if amount <= 0 || amount > 100 {
		return fmt.Errorf("can't decrement by %v as it is not between 0 - 100", amount)
	}

	total, _, err := database.GetTotal(c.FestivalCode)
	if err != nil {
		total = 0
	}

	if total < amount {
		amount = total
	}

	database.AddValue(-amount, c.FestivalCode)

	fmt.Println("-", amount)
	fmt.Println("New total of", total+amount)

	sendTotal(c, true)
	return nil
}
