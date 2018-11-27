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
	"strconv"
	"time"
	"os"
)


func GetEnv(key, defaultVal string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultVal
	}
	return key
}

var (
	host = GetEnv("users_host", "oilspill.ocf.berkeley.edu")
	port = GetEnv("users_post",  "5000")
	user = GetEnv("users_user", "postgres")
	password = GetEnv("users_password", "docker")
	dbname = GetEnv("users_dbname" , "historymap_users")

	jwtSecretKey = []byte(GetEnv("historymap_jwt_secret_key", "aVerySecretKey"))
)

var psqlInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

//Change to byte buffer maybe
type UserCred struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claim struct {
	userId int `json:"userId"`
	jwt.StandardClaims
}

const saltChars = "01234567890!@#$%^&*" +
	"abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

const saltCharLength = len(saltChars)

type MiddlewareAdapter func(http.Handler) http.Handler

type AuthHandler func(uidToken int, valid bool) http.Handler

func AuthMiddleware(authFunc AuthHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get claims
		var js interface{}
		err := json.NewDecoder(r.Body).Decode(&js)
		if err != nil {
			return
		}

		m := js.(map[string]interface{})
		tokenString := fmt.Sprint(m["token"])
		if tokenString != "" {
			uidToken, valid := authenticate(&tokenString)
			authHandler := authFunc(uidToken, valid)
			authHandler.ServeHTTP(w, r)
		}
	})
}

func authenticate(tokenString *string) (int, bool) {
	tkn, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		return -1, false
	}

	claim,ok := tkn.Claims.(jwt.MapClaims);

	if ok == false {
		log.Printf("Server Error: Problem Reading Claim")
		return -1, false
	}

	userId, err := strconv.Atoi(fmt.Sprint(claim["userId"]))
	if err != nil {
		return -1, false
	}

	if tkn.Valid {
		return userId, true
	} else {
		log.Printf("Invalid Jwt Token")
		return -1, false
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

func signup(w http.ResponseWriter, req *http.Request) {
	//check to see if username is unique
	x := new(UserCred)
	if err := json.NewDecoder(req.Body).Decode(x); err != nil {
		return; //server error
	}

	db, err := sql.Open("postgres", psqlInfo)
	defer db.Close()

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

	if checkName != "" {
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

type JwtToken struct {
	Token string `json:"token"`
}

func login(w http.ResponseWriter, req *http.Request) {
	var (
		expectedHash string
		salt         string
		id 			 int
	)

	x := new(UserCred)
	if err := json.NewDecoder(req.Body).Decode(x); err != nil {
		return
	}
	//include env variable for database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		return
	}

	//change to use env variable for accounts  table
	err = db.QueryRow("SELECT id, password_hash, salt FROM  accounts  WHERE "+
		"username=$1", x.Username).Scan(&id, &expectedHash, &salt)

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

	claims := Claim{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecretKey)
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

	authServ.HandleFunc("/login/", login)
	authServ.HandleFunc("/signup/", signup)
	//switch to TLS
	log.Fatal(authServ.ListenAndServe())
}




















