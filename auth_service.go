package historymap_auth

import (
	"crypto/hmac"
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pg"
	"golang.org/x/crypto/scrypt"
	"log"
	"net/http"
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

func AuthMiddleware() MiddlewareAdapter {
	return func(h http.Handler) http.Handler {
		return http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
			//get claims
			var js interface{}
			err := json.Unmarshal(r.Body, &js)
			m := f.(map[string]interface{})
			tokenString := js["token"]
			if tokenString != nil {
				auth := authenticate(&tokenString)
				if auth {
					h.ServeHTTP(r, w)
				} else {
					w.WriteHeader(http.StatusUnAuthorized)
				}
			} else {
				w.writeHeader(http.BadRequest)
			}
		})
	}
}

func authenticate(tokenString *string) (jwt.MapClaims, bool) {
	token, err = jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("aVerySecretKey")
	})
	if err != nil {
		return nil, false
	}

	claims, ok = token.Claims.(jwt.MapClaims)

	if ok != nil {
		log.Printf("Server Error: Failed to Extract Claim")
		return nil, false
	}

	if token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid Jwt Token")
		return nil, false
	}
}

func createSalt(saltLength int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	salt := make([]byte, saltLength)
	for i := range salt {
		salt[i] = saltChars[r1.Intn(saltCharLength)]
	}
	return string(salt)
}

func signup(w http.ResponseWriter, req *http.Request) int {
	//check to see if username is unique
	x := new(UserCred)
	if err = json.NewDecoder(req.Body).Decode(x); err != nil {
		return 0 //server error
	}

	db, err := sql.open("postgres", nil)
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

	if checkName != nil {
		w.WriteHeader(http.BadRequest)
		return
	}

	if len(x.Password) < 10 {
		w.WriteHeader(http.BadRequest)
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
	w.writeHeader(http.StatusOK)
}

func login(w http.ResponseWriter, req *http.Request) int {
	var (
		expectedHash string
		salt         string
	)

	x := new(UserCred)
	if err = json.NewDecoder(req.body).Decode(x); err != nil {
		return -5
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
	attempt, err = scrypt.Key([]byte(x.Password), []byte(salt), 32768, 8, 1, 32)

	if err != nil {
		log.Fatal(err)
		return
	}

	if attempt != x.expectedHash {
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
