package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	chats map[string]*Chat
	users map[string]*User

	broadcast chan *Message
	newChat   chan *Chat
}

func NewServer() *Server {
	server := Server{
		chats:     make(map[string]*Chat),
		users:     make(map[string]*User),
		broadcast: make(chan *Message),
		newChat:   make(chan *Chat),
	}
	go server.route()
	return &server
}

func (server *Server) getChat(chatId string) *Chat {
	return server.chats[chatId]
}

func (server *Server) getUser(userId string) (*User, error) {
	if user, ok := server.users[userId]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (server *Server) regNewUser(user *User) string {
	server.users[user.uid] = user
	return user.uid
}

func (server *Server) regNewChat(chat *Chat) string {
	server.chats[chat.uid] = chat
	server.newChat <- chat
	return chat.uid
}

func (server *Server) StartServer(ch chan struct{}) {
	defer close(ch)

	http.HandleFunc("/connection", server.connection)
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		newChatRequestHandler(server, w, r)
	})
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		newUserRequestHandler(server, w, r)
	})
	http.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		template(server, w, r)
	})
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		config(server, w, r)
	})

	log.Println("••• SERVER STARTED •••")
	http.ListenAndServe(":8070", nil) // Уводим utils сервер в горутину
}

func (server *Server) connection(w http.ResponseWriter, r *http.Request) {
	// авторизация
	user, err := server.getUser((r.URL.Query()["uid"])[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("...invalid uid")
		return
	}

	// ws соединение
	connection, _ := upgrader.Upgrade(w, r, nil)
	log.Printf("••• CONNECTION OPENED with user=%s •••\n", user.name)

	// сохранили
	user.Conn = connection // сохранили пользователю соединение
	user.ConnectionOpened = true
	//server.users[user.uid] = user // сохранили юзера в мэп
	user.writeChats()
	// отправили в горутины читать/писать
	go user.read()
	go user.write()
}

func (server *Server) route() {
	for {
		select {
		case msg := <-server.broadcast:
			log.Println("...message delivered to chat")
			msg.Chat.sendMessage(msg)
		case chat := <-server.newChat:
			log.Println("...newChat " + chat.uid)
		}

	}
}

//func (server *Server) newChat(holderId string, memberId string){
//	newRoom := newChat(server.users[holderId], server.users[memberId])
//	server.lasRoomId = server.lasRoomId + 1
//	server.chats[server.lasRoomId] = newRoom
//}

//func (server *Server) WriteMessage(message []byte) {
//	for _, user := range server.users {
//		err := user.Conn.WriteMessage(websocket.TextMessage, message)
//		if err != nil {
//			log.Println("write err: ", err)
//		}
//		log.Printf("-> msg: %s", string(message))
//	}
//}
