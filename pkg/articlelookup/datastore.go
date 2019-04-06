package articlelookup

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

/**Deals with storing new data that comes in from articlestore. **/
/**Get request from r and insert (or replace relevant field in db) **/
func store(db *sql.DB) {
	/**I believe that nchan will automatically deal with long polling**/
	for {
		resp, err := http.Get("nginx:6000/sub")
		if err != nil {
			log.Println(err)
			return
		}
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				continue
			}
			var articles []Marker
			err = json.Unmarshal(body, &articles)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, element := range body {
				/**Update Database**/
			}
		} else {
			log.Println(resp.StatusCode)
			continue
		}
	}
}
