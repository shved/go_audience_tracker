package main

import (
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
	videoIDs := customers[customerID]
	mutex.Unlock()

	set := uniqueSet(videoIDs)
	count := len(set)

	log.Println("customer stat called: ", customerID, count)
}

func videoCount(w http.ResponseWriter, r *http.Request) {
	mutex := &sync.Mutex{}
	videoID := parseIDFromURL(r.URL.Path)

	mutex.Lock()
	customerIDs := videos[videoID]
	mutex.Unlock()

	set := uniqueSet(customerIDs)
	count := len(set)

	log.Println("video stat called: ", videoID, count)
}

func uniqueSet(slice []int) []int {
	unique := make([]int, 0, len(slice))
	uniqueMap := make(map[int]bool)

	for _, val := range slice {
		if _, ok := uniqueMap[val]; !ok {
			uniqueMap[val] = true
			unique = append(unique, val)
		}
	}

	return unique
}

func parseIDFromURL(path string) (id int) {
	log.Println(path)
	stringSlice := strings.Split(path, "/")
	log.Println(stringSlice)
	id, _ = strconv.Atoi(stringSlice[2])
	log.Printf("ID %d, %T", id, id)
	return
}

// func sessionExpireTimeout(customerID, videoID int) {
// 	time.Sleep(6 * time.Second)
// }

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/pulse", pulse)
	mux.HandleFunc("/videos/", videoCount)
	mux.HandleFunc("/customers/", customerCount)

	http.ListenAndServe(":9292", mux)
}
