package articlestore

import (
	"database/sql"
	"net/http"
)

/** Adds user created articles/information to databases. Deals with making sure information is acceptable and keeps information stores **/

func customInfo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
