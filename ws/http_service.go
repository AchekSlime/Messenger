package ws

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserDto struct {
	name     string
	login    string
	password string
}

type ChatDto struct {
	users []string
}

func mapUser(user *UserDto, server *Server) *User {
	return NewUser(user.name, server)
}

func mapChat(chat *ChatDto, server *Server) *Chat {
	return NewChat(chat.users, server)
}

func newUserRequestHandler(server *Server, w http.ResponseWriter, r *http.Request) {
	var user UserDto
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("json mapper error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid := server.regNewUser(mapUser(&user, server))
	w.Write([]byte(uid))
}

func newChatRequestHandler(server *Server, w http.ResponseWriter, r *http.Request) {
	var chat ChatDto
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		log.Println("json mapper error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid := server.regNewChat(mapChat(&chat, server))
	w.Write([]byte(uid))
}

func config(server *Server, w http.ResponseWriter, r *http.Request) {
	userDto1 := UserDto{name: "achek"}
	userDto2 := UserDto{name: "egor"}

	uid1 := server.regNewUser(mapUser(&userDto1, server))
	uid2 := server.regNewUser(mapUser(&userDto2, server))

	chatDto := ChatDto{users: []string{uid1, uid2}}
	uidChat := server.regNewChat(mapChat(&chatDto, server))

	ans := uid1 + " " + uid2 + " " + uidChat

	w.Write([]byte(ans))
}
