package ws

import "messenger/models"

type Message struct {
	Uid      string
	Holder   *User
	Chat     *Chat
	Text     string
	SendTime int64
}

func (msg *Message) Wrap() *models.Message {
	return &models.Message{
		Uid:      msg.Uid,
		Holder:   msg.Holder.Wrap(),
		ChatId:   msg.Chat.uid,
		Text:     msg.Text,
		SendTime: msg.SendTime,
	}
}

func UnwrapMessage(protoMsg *models.Message, server *Server) *Message {
	user, _ := server.getUser(protoMsg.GetHolder().Uid)
	return &Message{
		Uid:      protoMsg.Uid,
		Holder:   user,
		Chat:     server.getChat(protoMsg.ChatId),
		Text:     protoMsg.Text,
		SendTime: protoMsg.SendTime,
	}
}
