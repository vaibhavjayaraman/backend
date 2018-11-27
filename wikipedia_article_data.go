package historymap_auth

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	"log"
	"net/http"
)

type articleRequest struct {
	Url string `json: "url"`
	Lat string `json: "lat"`
	Lon string `json: "lon"`
	Title string `json: "title"`
	ArticleInteraction string `json: "interactionType"`
}

var recordWikipediaData = AuthHandler(func(uidToken int, valid bool) http.Handler {
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

func recordData(request articleRequest) {
	db, err := gorm.Open("postgres", )
}


