package ws

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"net/http"
	"time"
)

type UserRegRequestDto struct {
	Name     string
	Login    string
	Password string
}

type ChatRequestDto struct {
	Users []string
}

type UserAuthRequestDto struct {
	Login    string
	Password string
}

type UserAuthResponseDto struct {
	Token string
	Uid   string
}

type Config struct {
	ChatType string
	ChatId   string
}

type Claims struct {
	jwt.StandardClaims
	Login string `json:"Login"`
}

var jwtKey = []byte("ergerh")

func mapUser(user *UserRegRequestDto, server *Server) *User {
	return NewUser(user.Name, user.Login, server)
}

func mapChat(chat *ChatRequestDto, server *Server) *Chat {
	return NewChat(chat.Users, server)
}

func regUserHandler(server *Server, w http.ResponseWriter, r *http.Request) {
	var user UserRegRequestDto
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("json mapper error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid := server.regNewUser(mapUser(&user, server))
	// toDo прокинуть статус
	w.Write([]byte(uid))
}

func authUser(server *Server, w http.ResponseWriter, r *http.Request) {
	var userDto UserAuthRequestDto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		log.Println("json mapper error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := server.findByLogin(userDto.Login)
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{ExpiresAt: jwt.At(time.Now().Add(time.Minute * 5))},
		Login:          user.login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("token err " + err.Error())
	}
	data, _ := json.Marshal(UserAuthResponseDto{Uid: user.uid, Token: tokenString})
	// toDo прокинуть статус
	w.Write(data)
	log.Println("...user <" + user.name + "> successfully authorized")
}

func regChatHandler(server *Server, w http.ResponseWriter, r *http.Request) {
	var chat ChatRequestDto
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

func parseToken(accessToken string) string {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected jwt auth method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return ""
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Login
	}

	return ""
}

func template(server *Server, w http.ResponseWriter, r *http.Request) {
	userDto1 := UserRegRequestDto{Name: "achek", Login: "achek"}
	userDto2 := UserRegRequestDto{Name: "egor", Login: "egor"}
	userDto3 := UserRegRequestDto{Name: "askar", Login: "askar"}
	userDto4 := UserRegRequestDto{Name: "semen", Login: "semen"}

	uid1 := server.regNewUser(mapUser(&userDto1, server))
	uid2 := server.regNewUser(mapUser(&userDto2, server))
	uid3 := server.regNewUser(mapUser(&userDto3, server))
	uid4 := server.regNewUser(mapUser(&userDto4, server))

	chatDto1 := ChatRequestDto{Users: []string{uid1, uid2}}
	uidChat1 := server.regNewChat(mapChat(&chatDto1, server))

	chatDto2 := ChatRequestDto{Users: []string{uid2, uid3, uid4}}
	uidChat2 := server.regNewChat(mapChat(&chatDto2, server))

	ans := userDto1.Login + " " + userDto2.Login + " " + userDto3.Login + " " + userDto4.Login + "\n" +
		uidChat1 + " " + uidChat2

	w.Write([]byte(ans))
}
