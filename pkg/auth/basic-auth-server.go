package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"historymap-microservices/pkg/tools"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)


var (
	host = tools.GetEnv("users_host", "oilspill.ocf.berkeley.edu")
	port = tools.GetEnv("users_post",  "5000")
	user = tools.GetEnv("users_user", "postgres")
	password = tools.GetEnv("users_password", "docker")
	dbname = tools.GetEnv("users_dbname" , "historymap_users")
)

var dbParams = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

type User struct {
	gorm.Model
	Name string `gorm:"size:255"`
	Username string `gorm:"size:255"`
	PasswordHash string `gorm:"type:text"`
	PasswordSalt string `gorm:"type:text"`
	Num int `gorm:"AUTO_INCREMENT"`
	Email  string  `gorm:"type:varchar(100);unique_index"`
	Joined time.Time
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func AuthServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login/", login)
	mux.HandleFunc("/signup/", signup)
}

func signup(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", dbParams)
	defer db.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := new(NewUser)
	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !stringsUnalike(newUser.Username, newUser.Name) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	login(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {

}

func stringsUnalike(a, b string) bool {
	return true
}

func createSalt(saltLength int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	salt := make([]byte, 245)
	for i := range salt {
		salt[i] = byte(strconv.FormatInt(int64(r1.Intn(255)), 2))
	}
	return string(salt)
}

