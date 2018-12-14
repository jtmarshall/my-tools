/*
Oauth2 Server that can serve as base authentication point for any app allowed access.
A seperate client server POSTs to the Goauth server with API Key in header and user/pass in form data.
If user/pass matches a 24hr token object is returned.

Setup:
	- Create an API Key to allow a POST to be made
	- Setup user database to store hashed passwords
	- Will need a client server to interact with this one
*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // Mysql driver
	"golang.org/x/crypto/bcrypt"
)

// OauthToken Schema
type OauthToken struct {
	User    string
	Token   string
	Role    int
	Expires time.Time
	YakPak  string
}

var (
	// API Key needed to accept request
	apiKey = "yourAPIKey"

	// DB connection creds
	host          = "hostURL"
	dbName        = "nameOfDatabase"
	user          = "userID"
	pass          = "%" + "userPassword"
	connectString = user + ":" + pass + "@tcp(" + host + ")/" + dbName
)

// Main
func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Index")
	})
	mux.HandleFunc("/login", handleLogin)

	// Heartbeat
	mux.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Ok")
		})

	// Start go routine for server to listen on port
	http.ListenAndServe(":"+port, mux)
}

// Handle Login POST request
func handleLogin(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/login" {
	// 	http.Error(w, "404 not found.", http.StatusNotFound)
	// 	return
	// }

	// Only process POST requests
	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			fmt.Fprintf(w, "false")
			return
		}

		// Check for API Key in header
		if r.Header.Get("X-Oauth-Key") != apiKey {
			fmt.Fprintf(w, "false")
			return
		}

		// Get user data from request
		user := r.FormValue("user")
		pass := r.FormValue("pass")

		// Get password hash from DB
		dbPass, dbRole, yakPak := getPwdRole(user)
		pwdMatch := checkPassword(dbPass, []byte(pass))

		// Write True if we validate user
		if pwdMatch {
			// Current time to byte string; pass this into hashing func to create unique token using bcrypt lib
			byteString := []byte(time.Now().String())

			newToken := &OauthToken{
				User:    user,
				Token:   hashSalt(byteString),
				Role:    dbRole,
				Expires: time.Now().Add(time.Minute * 2),
				YakPak:  yakPak,
			}

			// Convert Token to json so we can pass it back
			jsonToken, _ := json.Marshal(newToken)
			// Return json token
			fmt.Fprintf(w, string(jsonToken))
			return
		}

		// If no match return false
		fmt.Fprintf(w, "false")
		return
	default:
		fmt.Fprintf(w, "POST Only.")
	}
}

// Get user password and role
func getPwdRole(user string) (string, int, string) {
	// DB connect
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var dbPass string
	var role int
	var yakPak string

	// Select user's hashed password and role from DB
	row, err := db.Query("SELECT password, role, yakpak FROM user WHERE user = ?", user)
	if err != nil {
		log.Fatal(err)
	}
	for row.Next() {
		rowErr := row.Scan(&dbPass, &role, &yakPak)
		if rowErr != nil {
			log.Println("DB Password Retrieve Error:", rowErr)
		}
	}

	// Return user's hashed pass from DB as a byte slice
	return dbPass, role, yakPak
}

// Hash a password with golang library bcrypt
func hashSalt(pwd []byte) string {
	// GenerateFromPassword from bcrypt lib
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println("bycript error:", err)
	}

	// Return hashed password
	return string(hash)
}

// Check a password with golang library bcrypt
func checkPassword(hashedPass string, plainPass []byte) bool {
	// Convert hashed password from DB into byte slice
	byteHash := []byte(hashedPass)

	// Compare hash and password string with func from bcrypt lib
	err := bcrypt.CompareHashAndPassword(byteHash, plainPass)
	if err != nil {
		log.Println("Check Password Error:", err)
		return false
	}

	return true
}

// // Generate random string for tokens
// func randomString(n int) string {
// 	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letter[rand.Intn(len(letter))]
// 	}
// 	return string(b)
// }

// // Generate random numb in range to pass into randomString
// func randomNumb(min, max int) int {
// 	// Call Seed, using current nanoseconds to ensure numbers are random
// 	rand.Seed(int64(time.Now().Nanosecond()))
// 	return rand.Intn(max-min) + min
// }
