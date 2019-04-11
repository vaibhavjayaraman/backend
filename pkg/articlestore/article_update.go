package articlestore

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type latLonReq struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

/**Checks to make sure that incoming protobuffer information from articlelookup which includes lat/lon information is updated. Has another
method that sends updates via protobuffer to article_lookup for updating records **/

func articleStore() {
	connStr := ""
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/articles", articles(db))
	mux.HandleFunc("/interactions", interactions(db))
	mux.HandleFunc("/customInfo", customInfo(db))
}

func articles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			/*Write in Log */
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var llrq latLonReq
		err = json.Unmarshal(body, &llrq)
		if err != nil {
			/*Add to Log */
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		go processArticles(&llrq)
		w.WriteHeader(http.StatusOK)
	}
}

func processArticles(llrq *latLonReq) {

}

func sendArticles() {

}
