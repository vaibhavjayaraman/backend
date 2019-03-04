package tileserver

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var TILEROOT string = "/tiles/"
var regions []string = []string{"iberia", "mediaeval_middle_east", "northern_europe"}

func TileServer() {
	c := cache.New(48*60*time.Minute, 60*time.Minute)
	for i, region := range regions {
		go createRegion(c, region)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", findTile(c))
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func createRegion(c *cache.Cache, region string) {
	years , err := ioutil.ReadDir(TILEROOT + region)

	for i, year := range years
}

func GetTiles(c, region, year string) string {

}

func findTile(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		region := path[0]
		year := path[1]
		val := region + year
		resolvedYear, found := c.Get(val)
		if found {
			var rYear string = resolvedYear.(string)
			if rYear != "NA" {
				GetTiles(region, rYear)
			}
			w.WriteHeader(http.StatusOK)
			return
		} else {
			yr, err := strconv.Atoi(year)
			if err != nil {
				log.Fatal(err)
			}

			if yr > 1992 || yr < 0 {
				/**So that nobody puts a small or large number and makes us a ton of unnecessary dates **/
				w.WriteHeader(http.StatusOK)
				return
			}

			year = GetTiles(region, year)

			/**NA means that the time is before any time before **/
			if year == "NA" {
				w.WriteHeader(http.StatusOK)
				return
			}

			rYear, err := strconv.Atoi(year)
			if err != nil {
				log.Fatal(err)
			}

			for i := rYear; i <= yr; i++ {
				ryr := strconv.Itoa(i)
				val := region + ryr
				c.Set(val, year, cache.DefaultExpiration)
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
