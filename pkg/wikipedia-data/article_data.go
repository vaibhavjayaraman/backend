package wikipediadata

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"historymap-microservices/pkg/middleware"
	"historymap-microservices/pkg/tools"
	"log"
	"net/http"
	"time"
)

var (
	host = tools.GetEnv("wikipedia_data_host", "oilspill.ocf.berkeley.edu")
	port = tools.GetEnv("wikipedia_data_post",  "5000")
	user = tools.GetEnv("wikipedia_data_user", "postgres")
	password = tools.GetEnv("wikipedia_data_password", "docker")
	dbname = tools.GetEnv("wikipedia_data_dbname" , "historymap_wikipedia")
)

var dbParams = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

func ArticleDataServer() {
	db, err := gorm.Open("postgres", dbParams)
	defer db.Close()
	if err != nil {
		return
	}

	mux := http.NewServeMux()
	articles := make(chan articleData, 5000)
	users := make(chan userArticleData, 5000)

	go processArticleData(db, articles)
	go processUserData(db, users)

	authChain := middleware.Auth(recordData(false, articles, nil))
	mux.HandleFunc("/wikidata", authChain(recordData(true, articles, users)))
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func processArticleData (db *gorm.DB, in <-chan articleData) {
	for data :=  range in  {
		if db.NewRecord(data) {
			data.CreatedAt = time.Now()
			data.HoveredOver = 0
			data.Clicked = 0
			data.Generated = 0
			data.Searched = 0
		} else {
			db.First(&data)
		}

		switch data.ArticleInteraction {
		case GENERATED:
			data.Generated += 1
		case HOVERED_OVER:
			data.Generated += 1
		case CLICKED:
			data.Clicked += 1
		case SEARCHED:
			data.Searched += 1
		}
		data.UpdatedAt = time.Now()
		db.Save(data)
	}
}

func processUserData (db *gorm.DB, in <-chan userArticleData) {
	for data := range in {
		if db.NewRecord(data) {
			data.CreatedAt = time.Now()
			data.HoveredOver = 0
			data.Clicked = 0
			data.Generated = 0
			data.Searched = 0
		} else {
			db.First(&data)
		}

		switch data.ArticleInteraction {
		case GENERATED:
			data.Generated += 1
		case HOVERED_OVER:
			data.Generated += 1
		case CLICKED:
			data.Clicked += 1
		case SEARCHED:
			data.Searched += 1
		}
		data.UpdatedAt = time.Now()
		db.Save(data)
	}
}

func recordData(userAuth bool, articles chan<- articleData, users chan<- userArticleData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := new(articleRequest)
		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := articleData {
			Url: request.Url,
			Title: request.Title,
			Lat: request.Lat,
			Lon: request.Lon,
			ArticleInteraction: request.ArticleInteraction,
		}

		articles <- data

		if userAuth {
			userData := userArticleData{
				UserId: request.UserId,
				Url: request.Url,
				Title: request.Title,
				Lat: request.Lat,
				Lon: request.Lon,
				ArticleInteraction:request.ArticleInteraction,
			}

			users <- userData
		}
	}
}

type articleRequest struct {
	Url string `json: "url"`
	Lat float32 `json: "lat"`
	Lon float32 `json: "lon"`
	Title string `json: "title"`
	ArticleInteraction int `json: "interactionType"`
	UserId uint `json:name`
}

type userArticleData struct {
	gorm.Model
	UserId uint
	Url string
	Title string
	Lat float32
	Lon float32
	HoveredOver int
	Generated int
	Clicked int
	Searched int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	ArticleInteraction int `gorm:"-"`
}

type articleData struct {
	gorm.Model
	Url string
	Title string
	Lat float32
	Lon float32
	HoveredOver int
	Generated int
	Clicked int
	Searched int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	ArticleInteraction int `gorm:"-"`
}

