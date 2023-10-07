package lib

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func MakeRequest(method string, para string, url string, auth string) string {
	client := &http.Client{}
	data := strings.NewReader(para)
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		log.Fatal(err)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if auth != "" {
		//println(auth)
		req.Header.Set("Authorization", auth)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyText)
	return bodyString
}

func CheckToken(token string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:1323/api/info", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func CheckHTTPServer(url string) int { // 0: server is down, 1: server is up
	// 2: server bad gateway
	client := &http.Client{
		Timeout: 1 * time.Second, // 设置超时时间
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 1
	}
	return 2
}
