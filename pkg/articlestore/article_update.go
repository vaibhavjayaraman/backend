package articlestore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/tidwall/gjson"
	"github.com/worldhistorymap/articlelookup"
)

var WIKIPEDIA_PAGE_URL = "https://en.wikipedia.org/?curid="

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
	article_num 		 uint64
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
			/*Check what happens if non float Lat/Lon are added */
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
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

	
	[]articles := gjson.Get(body, "query.geosearch")

	if len(articles) >= 0 {
		for article := range articles {
			/**Enter into database **/
			/**Should this be on its own thread **/
			wikipediaUpdate(article)
		}
	}
}

func wikipediaUpdate(db *sql.DB, article []byte) {
	lat := gjson.Get(article, "lat");
	lat := gjson.Get(article, "lon");
	pgID := gjson.Get(article, "pageid");
	title := gjson.Get(article, "title");
	url := WIKIPEDIA_PAGE_URL + pgID
	info := getWikipediaExtract(title)
	mkr := marker {
		lat: lat,
		lon: lon, 
		url: url, 
		source: "wikipedia", 
		info: info, 
		title: title, 
	}
	updateDB(db, &mkr)
}

func getWikipediaExtract(title string) {
	url := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", title)
	resp, err := http.Get(url);
	if err != nil {
		log.Fatal(err)
		return
	}
	return gjson.get(resp.Body, "extract")
}

func dbqueue(mkrs chan<- * marker, mkr * marker) err {
	/**Put in channel to avoid concurrency issue when updating article num */
	mkrs <- mkr
}

func dbUpdate(mkrs <-chan * marker, db *sql.DB) {		
	for x := range mkrs {
		/**Call to DB and check if already in the DB**/
	}
}
func sendInfoArticleLookup(mkr * Marker) {	
	/**http.Post("nginx:6000/pub")**/
}

func sendInfoRecommendation(mkr *Marker) {
	return 
}

type ArticleNum struct {
	num uint64 	`json:"num"`
}

func getInfo(db *sql.DB) http.Handler {
	/**Query DB **/	
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

		var art ArticleNum
		err := json.Unmarshal(body, &art)		
		queryStmt, err := db.Prepare("SELECT url, info, title, source, lat, lon, num FROM markers WHERE num == ?;")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
		rows, err := queryStmt.Query(art.num)
		/**Check if article does not exist with that num**/
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

		if err := rows.Scan(&mkr) ; err != nil {
			/*Add Log */
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
		
		markerJson, err := json.Marshal(mkr)	
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
		w.Write(markerJson)
		w.WriteHeader(http.StatusOK)
	}
}