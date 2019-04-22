package articlestore

import (
	"database/sql"
	"net/http"
	"text/template"
)

/** Adds user created articles/information to databases. Deals with making sure information is acceptable and keeps information stores **/

type custom struct {
	Lat float64  `json:"lat"`
	Lon float64  `json:"lon"`
	Title string `json:"title"`
	Text string  `json:"text"`
	User string   `json:"user"`
}

func customInfo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.close()
		if err != nil {
			w.WriteHeader(http.StatusInternalError)
			return 
		}

		var csm custom
		err = json.Unmarshal(body, &csm)
		if err != nil {
			w.WriteHeader(http.StatusInternalError)
			return
		}	
		go addCustom(db &csm)
		w.WriteHeader(http.StatusOK)
	}
}

func addCustom(db * sql.DB, csm *custom) {
	var mkr marker
	marker.info = template.EscapeString(custom.Text)
	if len(marker.title) > 100 {
		/*Add to Log*/
		return 
	}
	marker.title = template.EscapeString(custom.Title)
	if len(marker.title) > 2000 {
		/*Add to Log */
		return 
	}
	marker.lat, err := strconv.ParseFloat(template.EscapeString(custom.Lat), 64)
	if err != nil {
		/*Add to Log*/
		return
	}
	marker.lon, err = strconv.ParseFloat(template.EscapeString(custom.Lon), 64)
	if err != nil {
		/*Add to Log */
		return
	}

	marker.source = fmt.Sprintf("User: %s", custom.User)
	if len(marker.source) > 100 {
		/*Add to Log */
		return 
	}

	err = updateDB(db, &marker)	
	if err != nil {
		/*Add to Log */
		return 
	}
}
