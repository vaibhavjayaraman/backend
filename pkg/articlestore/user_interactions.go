package articlestore

import (
	"database/sql"
	"net/http"
)

/**Stores information about how articles are interacted with. An api will also be created so that other services such as a recommendation service will
be able to easily query the information to make recommendations (kind of like how prometheus exporter uses metrics so that prometheus server can read information**/

func interactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
