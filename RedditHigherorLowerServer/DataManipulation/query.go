package datamanipulation

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
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
func getKey(byteArr []byte, key string) string {
	var output []byte
	var checkSequence string
	key = "\"" + key + "\""
	keyLen := len(key)
	for i := 0; i < int(len(byteArr)-keyLen); i++ {
		checkSequence = string(byteArr[i : i+keyLen])

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

func getRoutine(subredditData *SubredditData, wg *sync.WaitGroup) {
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

	//FIXME: Handle error states properly by changing to stored data

	request, requestErr = http.NewRequest("GET", "https://www.reddit.com/r/random/.json", nil)
	if requestErr != nil {
		log.Fatalln(requestErr)
	}
	request.Header.Set("User-Agent", "Custom Agent")
	targetSubreddit, _ = client.Do(request)

	remainingRequests, rateLimitError = strconv.ParseFloat(targetSubreddit.Header.Get("x-ratelimit-remaining"), 64)
	if rateLimitError != nil {
		log.Fatalln("Rate Limit encountered error,", rateLimitError)
	}

	if remainingRequests != 0.0 && targetSubreddit.StatusCode == 200 {
		bodyData, _ := io.ReadAll(targetSubreddit.Body)

		subredditName := getKey(bodyData, "subreddit_id")
		log.Printf("Successfuly queried random subreddit, id: %s. Querying subreddit data...\n", subredditName)

		request, requestErr = http.NewRequest("GET", "https://www.reddit.com/api/info.json?id="+subredditName, nil)
		if requestErr != nil {
			log.Fatal(requestErr)
		}
		request.Header.Set("User-Agent", "Custom Agent")
		targetSubreddit, _ = client.Do(request)
		bodyData, _ = io.ReadAll(targetSubreddit.Body)

		subredditData.Description = getKey(bodyData, "public_description")
		subredditData.Title = getKey(bodyData, "title")
		subredditData.Name = getKey(bodyData, "display_name")
		var atoiErr error
		subredditData.NumMembers, atoiErr = strconv.ParseInt(getKey(bodyData, "subscribers"), 10, 64)
		if atoiErr != nil {
			log.Fatal("Error with atoi, message:", atoiErr)
		}

		jsonMessage, _ := json.MarshalIndent(subredditData, " ", "    ")
		log.Println("Queried server data:", string(jsonMessage))

	} else {
		log.Println("Server rejected request, code: ", targetSubreddit.Status)
	}

	wg.Done()
}

// Get two random subreddits, will attempt to query Reddit, but if not possible will get a random subreddit stored serverside
// Will store queried subreddits for later use
func GetSubreddits() [2]SubredditData {
	start := time.Now()

	var subreddits [2]SubredditData
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)

	go getRoutine(&subreddits[0], waitGroup)
	go getRoutine(&subreddits[1], waitGroup)

	waitGroup.Wait()

	log.Printf("Time to get subreddit data: %dms\n", time.Since(start).Milliseconds())
	return subreddits
}
