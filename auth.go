package main

import (
	"crypto/hmac"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pg"
	"golang.org/x/crypto/scrypt"
	"log"
	"net/http"
)

//Change to byte buffer maybe
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, req *http.Request) int {
	var (
		expectedHash string
		salt         string
	)

	x = new(UserLogin)
	if err = json.NewCoder(req.body).Decode(x); err != nil {
		return -5
	}
	//include env variable for database
	db, err := sql.Open("postgres", nil)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	//change to use env variable for user_auth table
	err = db.QueryRow("SELECT password_hash, salt FROM  user_auth WHERE"+
		"username=$1", x.Username).Scan(&expectedHash, &salt)
	if err != nil {
		log.Fatal(err)
	}

	//check if password matches hashedpassword
	//change N, r, p to nonmagic num
	attempt, err = scrypt.Key([]byte(x.Password), []byte(salt), 32768, 8, 1, 32)

	if err != nil {
		log.Fatal(err)
	}

	if attempt != x.expectedHash {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//return jwt-go

}

func main() {

	//change address container env variable
	authServ := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	authServ.HandleFunc("/login", login)
	//switch to TLS
	log.Fatal(authServ.ListenAndServe())
}
