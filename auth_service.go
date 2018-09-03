package historymap_auth

import (
	"math/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/scrypt"
	"log"
	"net/http"
	"reflect"
	"time"
)

//change to use env var
var auth_db string = "database"

//Change to byte buffer maybe
type UserCred struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const saltChars = "01234567890!@#$%^&*" +
	"abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

const saltCharLength = len(saltChars)

type MiddlewareAdapter func(http.Handler) http.Handler

type AuthHandler func(username string) http.Handler

func AuthMiddleware(next AuthHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get claims
		var js interface{}
		err := json.NewDecoder(r.Body).Decode(&js);
		if err != nil {
			return;
		}
		m := js.(map[string]interface{})
		tokenString := fmt.Sprint(m["token"]);
		if tokenString != "" {
			username, auth := authenticate(&tokenString)
			if auth {
				authHandler := next(username);
				authHandler.ServeHTTP(w, r);
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	})
}

func authenticate(tokenString *string) (string, bool) {
	tkn, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("aVerySecretKey"), nil
	})
	if err != nil {
		return nil, false
	}

	claim,ok := tkn.Claims.(jwt.MapClaims);

	if ok == false {
		log.Printf("Server Error: Problem Reading Claim")
		return "", false
	}

	username := fmt.Sprint(claim["username"]);

	if tkn.Valid {
		return username, true
	} else {
		log.Printf("Invalid Jwt Token")
		return "", false
	}
}

func createSalt(saltLength int) string {
	s1 := rand.Seed(time.Now().UnixNano())
	r1 := rand.New(s1)
	salt := make([]byte, saltLength)
	for i := range salt {
		salt[i] = saltChars[r1.Intn(saltCharLength)]
	}
	return string(salt)
}

func signup(w http.ResponseWriter, req *http.Request) {
	//check to see if username is unique
	x := new(UserCred)
	if err := json.NewDecoder(req.Body).Decode(x); err != nil {
		return; //server error
	}

	db, err := sql.Open("postgres", nil)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return
	}

	var checkName string

	err = db.QueryRow("SELECT username FROM accounts WHERE "+
		"username=$1", x.Username).Scan(&checkName)
	if err != nil {
		log.Fatal(err)
		return
	}

	if checkName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(x.Password) < 10 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	salt := createSalt(20)
	encrypted, err := scrypt.Key([]byte(x.Password),
		[]byte(salt), 32768, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
		return
	}

	//create id INT PRIMARY KEY
	if _, err := db.Exec(
		"INSERT INTO accounts (username, password_hash, salt) "+
			"VALUES ($1, $2, $3)", x.Username, encrypted, salt); err != nil {
		log.Fatal(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func login(w http.ResponseWriter, req *http.Request) {
	var (
		expectedHash string
		salt         string
	)

	x := new(UserCred)
	if err := json.NewDecoder(req.Body).Decode(x); err != nil {
		return;
	}
	//include env variable for database
	db, err := sql.Open("postgres", nil)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		return
	}

	//change to use env variable for accounts  table
	err = db.QueryRow("SELECT password_hash, salt FROM  accounts  WHERE"+
		"username=$1", x.Username).Scan(&expectedHash, &salt)
	if err != nil {
		log.Fatal(err)
		return
	}

	//check if password matches hashedpassword
	//change N, r, p to nonmagic num
	attempt, err := scrypt.Key([]byte(x.Password), []byte(salt), 32768, 8, 1, 32)

	if err != nil {
		log.Fatal(err)
		return
	}

	if reflect.DeepEqual(string(attempt[:]), expectedHash) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claim := Claim{
		x.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			IssuedAt:  time.Now(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("aVerySecretKey"))
	if err != nil {
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	w.WriteHeader(http.StatusOK)
}

func main() {

	//change address container env variable
	authServ := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	authServ.HandleFunc("/login", login)
	authServ.HandleFunc("/signup", signup)
	//switch to TLS
	log.Fatal(authServ.ListenAndServe())
}
