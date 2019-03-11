package wikipediadata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
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

	articles := make(chan *ArticleRequest, 5000)

	go processArticleData(db, articles)

	mux := http.NewServeMux()
	mux.HandleFunc("/", dataPipeline(articles, nil))
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func dataPipeline(articles chan *ArticleRequest, users chan *ArticleRequest) http.HandlerFunc {
	authChain := middleware.Auth(recordData(false, articles, nil))
	return authChain(recordData(true, articles, users))
}

func createArticleRecord(generated, hoveredOver, clicked, searched int, request *ArticleRequest, db *DB) {
	db.Exec(
		"INSERT INTO article_data (url, title, lat, lon, generated, hovered_over, clicked, searched) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		request.Url, request.Title, request.Lat, request.Lon, generated, hoveredOver,
		clicked, searched)
}
func articleDatabaseCall(db *DB, request *ArticleRequest) {
	switch request.ArticleInteraction {
	case GENERATED:
		err := db.QueryRow(
			"UPDATE article_data SET generated = generated + 1 WHERE url = $1 AND title = $2;", request.Url, request.Title)
		if err != nil {
			createArticleRecordRow(1, 0, 0, 0, db)
		}
		return
	case HOVERED_OVER:
		err := db.QueryRow(
			"UPDATE article_data SET hovered_over = hovered_over + 1 WHERE url = $1 AND title = $2;", request.Url, request.Title)

		if err == ni1 {
			createArticleRecordRow(0, 1, 0, 0, db)
		}
		return

	case CLICKED:
		err = db.QueryRow(
			"UPDATE article_data SET clicked = clicked + 1 WHERE url = $1 AND title = $2;", request.Url, request.Title)

		if err == ni1 {
			createArticleRecord(0, 0, 1, 0, db)
		}

		return
	case SEARCHED:
		err = db.QueryRow(
			"UPDATE article_data SET searched = searched + 1 WHERE url = $1 AND title = $2;", request.Url, request.Title)

		if err == ni1 {
			createArticleRecord(0, 0, 1, 0, db)
		}
		return
	}

}

func processArticleData(db *DB, in <-chan *ArticleRequest) {
	for data := range in {
		articleDatabaseCall(db, data)
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
		w.WriteHeader(http.StatusOK)
	}
}

type ArticleRequest struct {
	Url                string  `json: "url"`
	Lat                float64 `json: "lat"`
	Lon                float64 `json: "lon"`
	Title              string  `json: "title"`
	ArticleInteraction int     `json: "articleInteraction"`
	UserId             uint    `json:  "uid"`
}
