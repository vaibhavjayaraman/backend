package articlelookup

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type LatLonReq struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func articleLookup() {
	connStr := ""
	db, err := sql.Open("postgres", connStr)
	mux := http.NewServeMux(db)
	mux.HandleFunc("/", handler(db))
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func handler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/**Use protobuff later **/
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var llrq LatLonReq
		err = json.Unmarshal(body, &llrq)

		/**Look for lat/lon from protobuff **/
		/** Do KNN search on underlying database.
		If nothing is returned within a certain distance, return false but return closest found if within some other larger range,
		and pass the lat/lon onto the articleservice api (as a go routine) so it can run the query to update the stores. If found, return true, and then
		add to the protobuff values **/
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
