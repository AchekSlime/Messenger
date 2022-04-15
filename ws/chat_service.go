package ws

import (
	"github.com/google/uuid"
	"log"
	"messenger/models"
)

type Chat struct {
	uid      string
	users    map[string]*User
	messages map[string]*Message
	chatType string

	server *Server
}

func NewChat(members []string, server *Server) *Chat {
	chat := Chat{
		uid:      uuid.New().String(),
		users:    make(map[string]*User),
		messages: make(map[string]*Message),
		chatType: "SINGLE",
		server:   server,
	}
	for _, v := range members {
		user, _ := server.getUser(v)
		chat.addNewUser(user)
	}
	return &chat
}

func (chat *Chat) Wrap() *models.Chat {
	members := make([]*models.User, 0)
	messages := make([]*models.Message, 0)
	for _, u := range chat.users {
		members = append(members, u.Wrap())
	}
	for _, m := range chat.messages {
		messages = append(messages, m.Wrap())
	}
	return &models.Chat{
		Uid:      chat.uid,
		ChatType: chat.chatType,
		Members:  members,
		Messages: messages,
	}
}

func (chat *Chat) sendMessage(msg *Message) {
	log.Println("...saving msg to chat")
	chat.messages[msg.Uid] = msg
	log.Println("...sending to other members")
	for _, user := range chat.users {
		if user.ConnectionOpened {
			user.broadcast <- msg
		}
	}
	log.Println("...has been send")
}

func (chat *Chat) addNewUser(user *User) {
	chat.users[user.uid] = user
	user.chats[chat.uid] = chat
	if len(chat.users) > 2 {
		chat.chatType = "MULTIPLE"
	}
}
