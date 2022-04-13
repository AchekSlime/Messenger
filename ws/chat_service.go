package ws

import "github.com/google/uuid"

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

func (room *Chat) sendMessage(msg *Message) {
	for _, user := range room.users {
		user.broadcast <- msg
	}
}
