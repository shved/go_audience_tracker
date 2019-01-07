package main

import (
	"github.com/deckarep/golang-set"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var customers = make(map[int][]int)
var videos = make(map[int][]int)

func pulse(w http.ResponseWriter, r *http.Request) {
	customerID, videoID := parseQuery(r.URL.RawQuery)
	mutex := &sync.Mutex{}

	mutex.Lock()
	customers[customerID] = append(customers[customerID], videoID)
	videos[customerID] = append(videos[videoID], customerID)
	mutex.Unlock()

	log.Printf("pulse with customer %d and video %d", customerID, videoID)
}

func parseQuery(rawQuery string) (customerID, videoID int) {
	values, _ := url.ParseQuery(rawQuery)
	customerID, _ = strconv.Atoi(values.Get("customer_id"))
	videoID, _ = strconv.Atoi(values.Get("video_id"))
	return
}

func customerCount(w http.ResponseWriter, r *http.Request) {
	mutex := &sync.Mutex{}
	customerID := parseIDFromURL(r.URL.Path)

	mutex.Lock()
	videoIDs := idsToInterfaceSlice(customers[customerID])
	mutex.Unlock()

	set := mapset.NewSetFromSlice(videoIDs)
	count := set.Cardinality()

	log.Println("customer stat called: ", count)
}

func videoCount(w http.ResponseWriter, r *http.Request) {
	mutex := &sync.Mutex{}
	videoID := parseIDFromURL(r.URL.Path)

	mutex.Lock()
	customerIDs := idsToInterfaceSlice(videos[videoID])
	mutex.Unlock()

	set := mapset.NewSetFromSlice(customerIDs)
	count := set.Cardinality()

	log.Println("video stat called: ", count)
}

func idsToInterfaceSlice(ids []int) []interface{} {
	s := make([]interface{}, len(ids))
	for i, v := range ids {
		s[i] = interface{}(v)
	}
	return s
}

func parseIDFromURL(path string) (id int) {
	stringSlice := strings.Split(path, "/")
	id, _ = strconv.Atoi(stringSlice[1])
	return
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/pulse", pulse)
	mux.HandleFunc("/videos/", videoCount)
	mux.HandleFunc("/customers/", customerCount)

	http.ListenAndServe(":9292", mux)
}
