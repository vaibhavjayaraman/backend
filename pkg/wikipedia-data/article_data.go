package wikipediadata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/github.com/lib/pq"
	"github.com/jinzhu/gorm"
	"github.com/worldhistorymap/backend/pkg/middleware"
	"github.com/worldhistorymap/backend/pkg/tools"
)

var (
	host     = tools.GetEnv("host", "oilspill.ocf.berkeley.edu")
	port     = tools.GetEnv("port", "5432")
	user     = tools.GetEnv("user", "postgres")
	password = tools.GetEnv("password", "docker")
	dbname   = tools.GetEnv("dbname", "historymap_wikipedia")
)

var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

func ArticleDataServer() {
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		return
	}

	articles := make(chan *request, 5000)
	users := make(chan *request, 5000)

	go processArticleData(db, articles)
	go processUserData(db, users)

	mux := http.NewServeMux()
	mux.HandleFunc("/", dataPipeline(articles, users))
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func dataPipeline(articles chan ArticleData, users chan UserArticleData) http.HandlerFunc {
	authChain := middleware.Auth(recordData(false, articles, nil))
	return authChain(recordData(true, articles, users))
}

var articleData = &new(ArticleData)
var userData = &new(UserArticleData)

func articleDatabaseCall(db *gorm.DB, request *ArticleRequest) {
	notFound := db.Where("url = ? AND title = ?", request.Url, request.Title).First(articleData)

	if notFound {
		articleData.Url = request.Url
		articleData.Title = request.Title
		articleData.CreatedAt = time.Now()
		articleData.Generated = 0
		aricleData.HoveredOver = 0
		articleData.Clicked = 0
		articleData.Searched = 0
	}

	switch request.ArticleInteraction {
	case GENERATED:
		articleData.Generated += 1
	case HOVERED_OVER:
		articleData.HoveredOver += 1
	case CLICKED:
		articleData.Clicked += 1
	case SEARCHED:
		articleData.Searched += 1
	}
	userData.UpdatedAt = time.Now()
	db.Save(articleData)
}

func userDatabaseCall(db *gorm.DB, request *ArticleRequest) {
	notFound := db.Where("url = ? AND title = ? AND user_id = ?", request.Url, request.Title, request.UserId).First(userData)

	if notFound {
		userData.Url = request.Url
		userData.Title = request.Title
		userData.UserId = request.UserId
		userData.CreatedAt = time.Now()
		userData.Generated = 0
		userData.HoveredOver = 0
		userData.Clicked = 0
		userData.Searched = 0
	}

	switch request.ArticleInteraction {
	case GENERATED:
		userData.Generated += 1
	case HOVERED_OVER:
		userData.HoveredOver += 1
	case CLICKED:
		userData.Clicked += 1
	case SEARCHED:
		userData.Searched += 1
	}
	userData.UpdatedAt = time.Now()
	db.Save(userData)
}

func processArticleData(db *gorm.DB, in <-chan *ArticleRequest) {
	for data := range in {
		articleDatabaseCall(db, data)
	}
}

func processUserData(db *gorm.DB, in <-chan *ArticleRequest) {
	for data := range in {
		userDatabaseCall(db, data)
	}
}

func recordData(userAuth bool, articles chan<- *ArticleRequest, users chan<- *ArticleRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := new(articleRequest)
		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		articles <- &request

		if userAuth {
			users <- &request
		}
		w.WriteHeader(http.StatusOK)
	}
}

type articleRequest struct {
	Url                string  `json: "url"`
	Lat                float64 `json: "lat"`
	Lon                float64 `json: "lon"`
	Title              string  `json: "title"`
	ArticleInteraction int     `json: "articleInteraction"`
	UserId             uint    `json:  "uid"`
}

type UserArticleData struct {
	gorm.Model
	UserId             uint   `gorm:"primary_key"`
	Url                string `gorm:"primary_key"`
	Title              string `gorm:"primary_key"`
	Lat                float64
	Lon                float64
	HoveredOver        int
	Generated          int
	Clicked            int
	Searched           int
	ArticleInteraction int `gorm:"-"`
}

type ArticleData struct {
	gorm.Model
	Url                string `gorm:"primary_key"`
	Title              string `gorm:"primary_key"`
	Lat                float64
	Lon                float64
	HoveredOver        int
	Generated          int
	Clicked            int
	Searched           int
	ArticleInteraction int `gorm:"-"`
}
