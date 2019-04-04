package articlelookup

import (
	"database/sql"
	"net/http"
)

/**Deals with storing new data that comes in from articlestore. **/
func store(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/**Get request from r and insert (or replace relevant field in db) **/
		return
	}
}
