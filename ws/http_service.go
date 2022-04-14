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

type Config struct {
	ChatType string
	ChatId   string
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

func template(server *Server, w http.ResponseWriter, r *http.Request) {
	userDto1 := UserDto{name: "achek"}
	userDto2 := UserDto{name: "egor"}

	uid1 := server.regNewUser(mapUser(&userDto1, server))
	uid2 := server.regNewUser(mapUser(&userDto2, server))

	chatDto := ChatDto{users: []string{uid1, uid2}}
	uidChat := server.regNewChat(mapChat(&chatDto, server))

	ans := uid1 + " " + uid2 + " " + uidChat

	w.Write([]byte(ans))
}

func config(server *Server, w http.ResponseWriter, r *http.Request) {
	user, err := server.getUser((r.URL.Query()["uid"])[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("...invalid uid")
		return
	}

	var chatUid string
	var config *Config
	if len(user.chats) == 1 {
		for _, v := range user.chats {
			chatUid = v.uid
		}
		config = &Config{"SINGLE", chatUid}
	} else {
		config = &Config{"MULTIPLE", ""}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(config)
	w.Write(data)
}
