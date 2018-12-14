package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// OauthToken Schema
type OauthToken struct {
	User    string
	Token   string
	Role    int
	Expires time.Time
}

var (
	oauthServerURL = "http://localhost:8080/login"          // URL for Oauth server
	apiKey         = "cc78b639-f142-4099-8f5d-2ff610f798c3" // API Key needed to for Oauth accept request
	activeUsers    = make(map[string]*OauthToken)           // Map of users/tokens to keep actives in memory, (make sure to use pointer to struct)
	redisClient    = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
)

// Run start up the router
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5050"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/oauth", oauthLogin)
	mux.HandleFunc("/updatePak", updatePak)

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)

	SetValue("foo", 0)
	SetValue("bar", 1)

	flushRedis()

	// Heartbeat
	mux.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Ok")
		})

	// Start go routine for server to listen on port
	http.ListenAndServe(":"+port, mux)
}

// Handle login attempt; Check Oauth server for valid user
func oauthLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fmt.Println("POST oauth")
		// Get user credentials
		user := r.FormValue("user")
		pass := r.FormValue("pass")

		// Check active user; Evaluates to ok if user(key) in map
		if val, ok := activeUsers[user]; ok {
			// Check if expire time is not up, return current user token
			if val.Expires.After(time.Now()) {
				userJSON, _ := json.Marshal(activeUsers[user])
				fmt.Fprintf(w, string(userJSON))
				//fmt.Println(activeUsers[user])
				return
			}
		}

		// If user not already in map requestAuth and if success user will be inserted into map, return user token
		if requestAuth(user, pass) {
			userJSON, _ := json.Marshal(activeUsers[user])
			fmt.Fprintf(w, string(userJSON))
			return
		}

		// If user not active AND can't be validated
		fmt.Fprintf(w, "Access Denied")
		return
	default:
		fmt.Fprintf(w, "POST Only.")
	}
}

// Request a new user token
func requestAuth(user string, pass string) bool {
	// Setup Form data for request
	form := url.Values{}
	form.Add("user", user)
	form.Add("pass", pass)

	// Setup request with our form data, then set headers for API Key and form
	req, err := http.NewRequest("POST", oauthServerURL, strings.NewReader(form.Encode()))
	req.Header.Set("X-Oauth-Key", apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// Get response body from oauth server
	body, _ := ioutil.ReadAll(resp.Body)

	// If false resp from goauth: return
	if string(body) == "false" {
		log.Println("Invalid User/Pass")
		return false
	}
	// If response doesn't return false update active user map
	if string(body) != "false" {
		// Initialize new token struct to hold resp data
		var respToken OauthToken

		// Decode JSON response and assign it to respToken
		jsonErr := json.Unmarshal(body, &respToken)
		if jsonErr != nil {
			log.Println("JSON encode ERR:", jsonErr)
			return false
		}

		// Save authenticated user in map
		activeUsers[user] = &respToken
	}
	return true
}

// Handle post call to update user's yakpak
func updatePak(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		user := r.FormValue("user")

		// Check active user; Evaluates to ok if user(key) in map
		if _, ok := activeUsers[user]; ok {
			reqPak := r.FormValue("yakpak")
			userPak, _ := GetValue(user)

			// If the pak sent with the request is not the same as the userPak we have in redis,
			// then update redis with the pak sent in the request
			if reqPak != userPak {
				SetValue(user, reqPak)
			}
			fmt.Fprint(w, "pak update")
		}

	default:
		fmt.Fprint(w, "failed pak update")
	}
}

// Set value in redis
func SetValue(key string, value interface{}) (bool, error) {
	serializedValue, _ := json.Marshal(value)
	err := redisClient.Set(key, string(serializedValue), 0).Err()
	return true, err
}

// Get value from redis
func GetValue(key string) (interface{}, error) {
	var deserializedValue interface{}
	serializedValue, err := redisClient.Get(key).Result()
	json.Unmarshal([]byte(serializedValue), &deserializedValue)
	return deserializedValue, err
}

// Flush redis store to db
func flushRedis() {
	var cursor uint64
	var n int
	// Loop and scan through redis keys inserting into 'keys' array
	for {
		var keys []string
		var err error
		// Scan all keys: (match '*')
		keys, cursor, err = redisClient.Scan(cursor, "*", 10).Result()
		if err != nil {
			log.Println(err)
			return
		}
		n += len(keys)

		fmt.Println(keys)

		if cursor == 0 {
			break
		}
	}

	fmt.Println("flushed")
}