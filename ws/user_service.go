package ws

import (
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"messenger/models"
)

type User struct {
	uid  string
	name string

	server           *Server
	Conn             *websocket.Conn
	ConnectionOpened bool

	notReadMsg []*Message
	broadcast  chan *Message
}

func NewUser(name string, server *Server) *User {
	return &User{
		uid:       uuid.New().String(),
		name:      name,
		broadcast: make(chan *Message),
		server:    server,
	}
}

func (user *User) Wrap() *models.User {
	return &models.User{
		Uid:  user.uid,
		Name: user.name,
	}
}

func (user *User) handleNewMessage(protoMsg *models.Message) {
	msg := UnwrapMessage(protoMsg, user.server)
	log.Println("...message unwrapped")
	user.server.broadcast <- msg
	log.Println("...message pushed into chat")
}

func (user *User) read() {
	defer func() {
		// todo трекаем что конекшн разорвался
		// todo закрываем конекшн
	}()

	for {
		_, data, err := user.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				user.ConnectionOpened = false
				log.Printf("error: %v", err)
			}
			break
		}
		log.Println("...bytes have been read")

		newMessage := new(models.Message)
		err = proto.Unmarshal(data, newMessage)
		log.Println("...proto unmarshaled")
		user.handleNewMessage(newMessage)
	}
}

func (user *User) write() {
	for {
		select {
		case msg := <-user.broadcast:
			protoMsg := msg.Wrap()
			data, _ := proto.Marshal(protoMsg)
			user.Conn.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}

func (user *User) writeNotReadMsg() {
	if len(user.notReadMsg) < 1 {
		return
	}

	messages := make([]*models.Message, 5)
	for _, msg := range user.notReadMsg {
		messages = append(messages, msg.Wrap())
	}

	protoNewMessages := &models.NewMessages{Messages: messages}
	data, _ := proto.Marshal(protoNewMessages)
	user.Conn.WriteMessage(websocket.BinaryMessage, data)
}
