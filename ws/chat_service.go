package ws

import (
	"github.com/google/uuid"
	"log"
)

type Chat struct {
	uid   string
	users map[string]*User

	server *Server
}

func NewChat(members []string, server *Server) *Chat {
	chat := Chat{
		uid:    uuid.New().String(),
		users:  make(map[string]*User),
		server: server,
	}
	for _, v := range members {
		user, _ := server.getUser(v)
		chat.users[v] = user
	}
	return &chat
}

func (chat *Chat) sendMessage(msg *Message) {
	log.Println("..sending to other members")
	for _, user := range chat.users {
		user.broadcast <- msg
	}
	log.Println("...has been send")
}
