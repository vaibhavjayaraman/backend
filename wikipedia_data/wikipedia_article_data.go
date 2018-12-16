package wikipedia_data


import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	"log"
	"net/http"
	"historymap-microservices/historymap_auth"
	"time"
)

var (
	host = historymap_auth.GetEnv("wikipedia_data_host", "oilspill.ocf.berkeley.edu")
	port = historymap_auth.GetEnv("wikipedia_data_post",  "5000")
	user = historymap_auth.GetEnv("wikipedia_data_user", "postgres")
	password = historymap_auth.GetEnv("wikipedia_data_password", "docker")
	dbname = historymap_auth.GetEnv("wikipedia_data_dbname" , "historymap_wikipedia")
)

var wikipediaDatabase = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

type articleRequest struct {
	Url string `json: "url"`
	Lat float32 `json: "lat"`
	Lon float32 `json: "lon"`
	Title string `json: "title"`
	ArticleInteraction string `json: "interactionType"`
}

type UserWikipediaData struct {
	gorm.Model
	UserId int
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
}

type WikipediaData struct {
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
}

var recordWikipediaData = historymap_auth.AuthHandler(func(uidToken int, valid bool) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		request := articleRequest{}
		log.Println(r.Method)

		requestJson, err := ioutil.ReadAll(r.Body);
		if err != nil {
			log.Fatal("Error reading from body", err)
		}


		err = json.Unmarshal(requestJson, &request)
		if err != nil {
			log.Fatal("Error Unmarshaling json", err)
		}

	})
})

func recordUserData(request articleRequest) {
	db, err := gorm.Open("postgres", wikipediaDatabase)
	defer db.Close()

	if err != nil {
		log.Fatal("GORM unable to connect to wikipediaDatabase", err)
	}


	var user = UserWikipediaData{}
	user.Url = request.Url
	user.Title = request.Title
	user.Lat = request.Lat
	user.Lon = request.Lon

	entryExists := db.NewRecord(&user)
	if !entryExists {
		db.Create(&user)
	} }


