package articlestore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type latLonReq struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type marker struct {
	url          string
	info         string
	title        string
	lat          float64
	lon          float64
	source       string
	generated    int64
	beg_year     int32
	end_year     int32
	hovered_over int64
	clicked      int64
	searched     int64
	created_at   time.Time
	updated_at   time.Time
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
		go findInfo(&llrq, db)
		w.WriteHeader(http.StatusOK)
	}
}

func findInfo(llrq *latLonReq, db *sql.DB) {
	findWikipedia(llrq, db)
}

func findWikipedia(llrq *latLonReq, db *sql.DB) {
	wikiRange := 9999
	fileReturnLimit := 10
	url := fmt.Sprintf("https://en.wikipedia.org/w/api.php?"+
		"action=query&origin=*&list=geosearch&gscoord=%f|%f"+
		"&gsradius=%d&gslimit=%d&prop=info|extracts&inprop=url"+
		"&format=json", llrq.Lat, llrq.Lon, wikiRange, fileReturnLimit)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	body := resp.Body
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
		return
	}

	var articles []map[string]map[string]interface{} = data["query"]["geosearch"]
	if len(articles) >= 0 {
		for article := range articles {
			/**Enter into database **/
			wikipediaUpdate()
		}
	}
}

func wikipediaUpdate(db *sql.DB /**Mapping for article**/) {
	var mkr Marker
	updateDB(db, mkr)
}
func updateDB(db *sql.DB, mkr *marker) {

}
func sendInfo() {

}
