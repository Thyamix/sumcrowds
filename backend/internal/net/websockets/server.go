package websockets

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/thyamix/festival-counter/internal/apperrors"
	"github.com/thyamix/festival-counter/internal/database"
)

type Server struct {
	Clients map[*Client]bool

	Messages chan Message

	Register chan *Client

	Unregister chan *Client
}

type Message struct {
	Message      []byte
	FestivalCode string
}

var server *Server

func NewHub() *Server {
	server = &Server{
		Clients:    make(map[*Client]bool),
		Messages:   make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	return server
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.Register:
			s.Clients[client] = true
		case client := <-s.Unregister:
			if ok := s.Clients[client]; ok {
				delete(s.Clients, client)
				close(client.Send)
			}
		case msg := <-s.Messages:
			for client := range s.Clients {
				if client.FestivalCode == msg.FestivalCode {
					select {
					case client.Send <- msg.Message:
					default:
						close(client.Send)
						delete(s.Clients, client)
					}
				}
			}
		}
	}

}

func BroadcastTotal(festivalCode string) error {
	total, maxGauge, err := database.GetTotalAndMax(festivalCode)
	if err != nil {
		log.Print(err)
		return apperrors.ErrFailedGetTotal
	}
	totalJson, err := json.Marshal(map[string]int{
		"total": total,
		"jauge": maxGauge,
	})

	if err != nil {
		return apperrors.ErrFailedMarshal
	}

	server.Messages <- Message{Message: totalJson, FestivalCode: festivalCode}

	fmt.Printf("Sending total:max:festivalCode: %v:%v:%v \n", total, maxGauge, festivalCode)

	return nil
}
