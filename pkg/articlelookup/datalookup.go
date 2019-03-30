package articlelookup

import (
	"log"
	"net/http"
)

func articleLookup() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func handler(w http.ResponseWriter, r *http.Request) {

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
- if cache miss (cannot find anything within 1000 kilometers (but work on this heuristic including a null marker) call wikipedia api yourself.
	- When passing off lat/lon via gRPC also note that

kafka - keeps master copy of article database
- once updated, then sends information via grpc to store in cached postgis database

-another golang service (datainteraction)
	- have content
	- have /error handle errors
	- have /interaction handle likes, times generated, times hovered over, and times clicked on - pass off to message bus

- golang service for recommendations (recommendations)
	- has own database with user recommendations - waits on pushes from central database
**/

/**
gRPC
	1)
		lat/lon/database miss
**/
