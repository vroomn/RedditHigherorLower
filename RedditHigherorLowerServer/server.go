package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	response, err := http.Get("https://www.reddit.com/reddits.json?count=1")
	if err != nil {
		fmt.Println(err)
		return
	}

	headerBytes, err := json.MarshalIndent(response.Header, " ", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(headerBytes))
	fmt.Println(response.StatusCode)

}
