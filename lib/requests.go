package lib

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type Code struct {
	Code string `json:"code"`
}

func SignUp(username string) (password string) {
	println("username=" + username)
	text := MakeRequest("POST", "username="+username, "http://127.0.0.1:1323/signup", "")
	//println("[Signup]" + text)
	println("[Signup]Success")
	var user User
	err := json.Unmarshal([]byte(text), &user)
	if err != nil {
		fmt.Println("[SignUp]Error decoding JSON:", err)
		return
	}
	return user.Password
}

func Login(username string, password string) (tokens string, errors error) {
	text := MakeRequest("POST", "username="+username+"&password="+password, "http://127.0.0.1:1323/login", "")
	//println("[Login]" + text)
	println("[Login]Success")
	var token Token
	err := json.Unmarshal([]byte(text), &token)
	if err != nil {
		fmt.Println("[Login]Error decoding JSON:", err)
		return
	}
	if token.Token == "" {
		return "", fmt.Errorf("Login failed")
	}
	return token.Token, nil
}

func HeartBeat(token string) (returnToken string) {
	text := MakeRequest("GET", "", "http://127.0.0.1:1323/api/heartbeat", "Bearer "+token)
	//println("[HeartBeat]" + text)
	println("[HeartBeat]Successfully refreshed token")
	var newToken Token
	err := json.Unmarshal([]byte(text), &newToken)
	if err != nil {
		fmt.Println("[HeartBeat]Error decoding JSON:", err)
		return
	}
	return newToken.Token
}

func Info(token string) (returnCode string) {
	text := MakeRequest("GET", "", "http://127.0.0.1:1323/api/info", "Bearer "+token)
	//println("[Info]" + text)
	println("[Info]Get code successfully")
	var code Code
	err := json.Unmarshal([]byte(text), &code)
	if err != nil {
		fmt.Println("[Info]Error decoding JSON:", err)
		return
	}
	println("[Info]Code: " + code.Code)
	return code.Code
}

func Validate(token string, code string) string {
	text := MakeRequest("POST", "code="+code, "http://127.0.0.1:1323/api/validate", "Bearer "+token)
	print("[Validate]")
	return text
}
