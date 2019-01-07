package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var customers []uint32
var videos []uint32

func pulse(w http.ResponseWriter, r *http.Request) {
	customerID, videoID := parseQuery(r.URL.RawQuery)
	newSession := session{VideoID: videoID, CustomerID: customerID}
	mutex := &sync.Mutex{}

	mutex.Lock()
	customers[customerID] = append(customers[customerID], videoID)
	videos[customerID] = append(videos[vidoeID], customerID)
	mutex.Unlock()

	log.Printf("pulse with customer %d and video %d", newSession.CustomerID, newSession.VideoID)
}

func parseQuery(rawQuery string) (customerID, videoID uint64) {
	values, _ := url.ParseQuery(rawQuery)
	customerID, _ = strconv.ParseInt(values.Get("customer_id"))
	videoID, _ = strconv.ParseInt(values.Get("video_id"))
	return
}

func videoStat(w http.ResponseWriter, r *http.Request) {
	mutex := &sync.Mutex{}
	videoID := parseIDFromURL(r.URL.Path)

	mutex.Lock()
	count := len(videos[videoID])
	mutex.Unlock()

	log.Println("video stat called: ", count)
}

func customerStat(w http.ResponseWriter, r *http.Request) {
	mutex := &sync.Mutex{}
	customerID := parseIDFromURL(r.URL.Path)

	mutex.Lock()
	count := len(customers[customerID])
	mutex.Unlock()

	log.Println("customer stat called: ", count)
}

func parseIDFromURL(path string) (id uint64) {
	stringSlice := strings.Split(path, "/")
	id, _ = strconv.ParseInt(stringSlice[1])
	return
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/pulse", pulse)
	mux.HandleFunc("/videos/", videoStat)
	mux.HandleFunc("/customers/", customerStat)

	http.ListenAndServe(":9292", mux)
}
