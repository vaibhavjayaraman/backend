package articlelookup

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type LatLonReq struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Year int     `json:"year"`
}

type Marker struct {
	url    string
	info   string
	title  string
	source string
	lat    float64
	lon    float64
}

func articleLookup() {
	connStr := ""
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler(db))
	go store(db)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func handler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/**Use protobuff later **/
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			/*Add in log */
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var llrq LatLonReq
		err = json.Unmarshal(body, &llrq)
		if err != nil {
			/*Add in log*/
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		/**Optimize this query **/
		/** KNN search on underlying database **/
		queryStmt, err := db.Prepare("SELECT url, info, title, source, lat, lon FROM markers WHERE beg_year <= $1 AND end_year >= $1" +
			"ORDER BY geom <-> ST_SetSRID(ST_MakePoint($2, $3), 4326) LIMIT 10;")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var markers []Marker
		rows, err := queryStmt.Query(llrq.Year, llrq.Lon, llrq.Lat)
		defer rows.Close()

		for rows.Next() {
			var marker Marker
			if err := rows.Scan(&marker); err != nil {
				/*Add in log */
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			markers = append(markers, marker)
		}
		markerJson, err := json.Marshal(markers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		go forwardCoords(body)
		w.Write(markerJson)
		w.WriteHeader(http.StatusOK)
	}
}

func forwardCoords(msg []byte) {
	/**switch to grpc**/
	_, err := http.Post("http://articlestore/markers", "json", bytes.NewReader(msg))
	if err != nil {
		/**Add in Logging **/
		return
	}
}

//create gRPC for lat/lon and have postgres instance for database be named something else
//https://docs.microsoft.com/en-us/sql/relational-databases/spatial/query-spatial-data-for-nearest-neighbor?view=sql-server-2017

/**
have / handle lat/lon
		- pass off lat/lon to message bus so that "trending" can be generated via mypy microservice
		- store lat/lon points generated at a later time (too much data otherwise)
-have /custommarker will add custom marker to the datastore
- golang database should be stupid. Just has the url, title, and some extract which it will continue to serve
- will also have a routine that looks to update underlying postgres via messagebus (and grpc)
- if cache miss (cannot find anything within 1000 kilometers (but work on this heuristic including a null marker) call
wikipedia api yourself. - When passing off lat/lon via gRPC also note that kafka - keeps master copy of article database -
once updated, then sends information via grpc to store in cached postgis database -another golang service (datainteraction) -
have content - have /error handle errors - have /interaction handle likes, times generated, times hovered over, and times clicked on - pass off to message bus
- golang service for recommendations (recommendations)
	- has own database with user recommendations - waits on pushes from central database
**/

/**
gRPC
	1)
		lat/lon/database miss
**/
