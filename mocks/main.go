package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your_secret_key")

// UserCredentials - структура для входных данных
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWTResponse - структура для ответа с JWT токеном
type JWTResponse struct {
	Token string `json:"token"`
}

// NodeInfo - структура для mock-ответа о нодах
type NodeInfo struct {
	NodeID     string `json:"nodeId"`
	UsersCount int    `json:"users_count"`
	UsersLimit int    `json:"users_limit"`
	State      string `json:"state"`
}

// GenerateJWT - функция для генерации JWT токена
func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// LoginHandler - обработчик для маршрута /api/user/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверка логина и пароля
	if creds.Username == "admin" && creds.Password == "admin" {
		token, err := GenerateJWT(creds.Username)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		// Отправляем токен в ответе
		response := JWTResponse{Token: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// NodeListHandler - mock-обработчик для маршрута /api/user/list/nodes
func NodeListHandler(w http.ResponseWriter, r *http.Request) {
	// Пример данных о нодах
	nodes := []NodeInfo{
		{NodeID: "some_node_name", UsersCount: 1, UsersLimit: 10, State: "online"},
		{NodeID: "foo", UsersCount: 3, UsersLimit: 5, State: "offline"},
	}

	// Отправка JSON-ответа с информацией о нодах
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func AvatarHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./default.png")
}

func main() {
	http.HandleFunc("/api/user/login", LoginHandler)
	http.HandleFunc("/api/user/list/nodes", NodeListHandler)
	http.HandleFunc("/api/user/avatar/default.png", AvatarHandler)
	http.ListenAndServe(":8090", nil)
}
