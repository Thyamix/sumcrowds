package websockets

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

func NewHub() *Server {
	return &Server{
		Clients:    make(map[*Client]bool),
		Messages:   make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
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
