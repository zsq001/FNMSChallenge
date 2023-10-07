package main

import (
	"FNMSChallenge/lib"
	"fmt"
	"log"
	"time"
)

// check server status then try to register
func register(username string) string {
	for true {
		if lib.CheckHTTPServer("http://127.0.0.1:1323/") == 1 {
			return lib.SignUp(username)
		}
		println("[Register]Server is down or bad gateway")
		println("[Register]Re-try in 10s")
		time.Sleep(10 * time.Second)
	}
	return ""
}

// check server status then try to log in
func login(username string, password string) (token string, pass string) {
	for true {
		newpass := password
		if lib.CheckHTTPServer("http://127.0.0.1:1323/") == 1 {
			token, err := lib.Login(username, newpass)
			if err != nil {
				newpass = register(username)
				token, err = lib.Login(username, newpass)
			}
			return token, newpass
		}
		println("[Login]Server is down or bad gateway")
		println("[Login]Re-try in 10s")
		time.Sleep(10 * time.Second)
	}
	return "", ""
}

// every 15 seconds check token and heartbeat then submit code
func main() {
	var username string
	fmt.Print("Enter a username: ")
	_, err := fmt.Scanln(&username)
	if err != nil {
		log.Fatalf("Error reading input: %s", err)
	}
	var status int
	var password, token string
	var count int64
	count = 0
	status = 0 //0->down 200->up 502->bad gateway
	password = register(username)
	//502->re-login shutdown->re-register
	for true {
		count++
		if lib.CheckHTTPServer("http://127.0.0.1:1323/") != 1 {
			if lib.CheckHTTPServer("http://127.0.0.1:1323/") == 0 {
				println("Server is down, re-try in 10s")
				status = 0
				time.Sleep(10 * time.Second)
				continue
			} else {
				println("Server is bad gateway!, re-try in 10s")
				status = 502
				time.Sleep(10 * time.Second)
				continue
			}
		}
		if status == 0 {
			password = register(username)
			token, password = login(username, password)
			count = 1
			status = 200
		}
		if status == 502 {
			token, password = login(username, password)
			count = 1
			status = 200
		}
		if token == "" || lib.CheckToken(token) == false {
			token, password = login(username, password)
			if token == "" {
				println("Login failed,Server maybe internal error,re-try in 10s")
				status = 502
				count = 1
				time.Sleep(10 * time.Second)
				continue
			}
			println("Login success!")
			count = 1
		}
		if (count % 15) == 1 {
			token = lib.HeartBeat(token)
			code := lib.Info(token)
			fmt.Println(lib.Validate(token, code))
		}
		println("[Info]Server status:" + fmt.Sprintf("%d", status))
		time.Sleep(1 * time.Second)
	}
}
