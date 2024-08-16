package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type SubredditData struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	NumMembers  int64  `json:"numMembers"`
}

/*
Looks for case-specific match to key in the byte array, automatically adds quotations before and after key

Intended to work with JSON type strings
*/
func GetKey(byteArr []byte, key string) string {
	var output []byte
	var checkSequence string
	key = "\"" + key + "\""
	keyLen := len(key)
	for i := 0; i < int(len(byteArr)-keyLen); i++ {
		checkSequence = string(byteArr[i : i+keyLen])

		//fmt.Println(checkSequence)
		if checkSequence == key {
			i += keyLen + 3
			for byteArr[i] != '"' && byteArr[i] != ',' {
				output = append(output, byteArr[i])
				i++
			}
			break
		}
	}

	return string(output)
}

// Get two random subreddits, will attempt to query Reddit, but if not possible will get a random subreddit stored serverside
// Will store queried subreddits for later use
func GetSubreddits() [2]SubredditData {
	var subreddits [2]SubredditData

	var targetSubreddit *http.Response
	var client = http.Client{
		Transport: &http.Transport{
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		},
	}
	var requestErr error
	var request *http.Request

	var remainingRequests float64
	var rateLimitError error

	for i := 0; i < 2; i++ {
		//FIXME: Handle error states properly by changing to stored data

		request, requestErr = http.NewRequest("GET", "https://www.reddit.com/r/random/.json", nil)
		if requestErr != nil {
			log.Fatal(requestErr)
		}
		request.Header.Set("User-Agent", "Custom Agent")
		targetSubreddit, _ = client.Do(request)

		remainingRequests, rateLimitError = strconv.ParseFloat(targetSubreddit.Header.Get("x-ratelimit-remaining"), 64)
		if rateLimitError != nil {
			fmt.Println(rateLimitError)
			log.Fatalln("Rate Limit encountered error")
		}

		if remainingRequests != 0.0 && targetSubreddit.StatusCode == 200 {
			bodyData, _ := io.ReadAll(targetSubreddit.Body)

			subredditName := GetKey(bodyData, "subreddit_id")
			log.Printf("Successfuly queried random subreddit, id: %s. Querying subreddit data...\n", subredditName)

			request, requestErr = http.NewRequest("GET", "https://www.reddit.com/api/info.json?id="+subredditName, nil)
			if requestErr != nil {
				log.Fatal(requestErr)
			}
			request.Header.Set("User-Agent", "Custom Agent")
			targetSubreddit, _ = client.Do(request)
			bodyData, _ = io.ReadAll(targetSubreddit.Body)

			subreddits[i].Description = GetKey(bodyData, "public_description")
			subreddits[i].Title = GetKey(bodyData, "title")
			subreddits[i].Name = GetKey(bodyData, "display_name")
			var atoiErr error
			subreddits[i].NumMembers, atoiErr = strconv.ParseInt(GetKey(bodyData, "subscribers"), 10, 64)
			if atoiErr != nil {
				fmt.Println("Error with atoi, message:", atoiErr)
			}

			jsonMessage, _ := json.MarshalIndent(subreddits[i], " ", "    ")
			log.Println("Queried server data:", string(jsonMessage))

		} else {
			fmt.Println("Server rejected request, code: ", targetSubreddit.Status)
		}

	}

	return subreddits
}

func main() {
	GetSubreddits()
}
