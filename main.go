package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var customers = make(map[int]map[int]int)
var videos = make(map[int]map[int]int)

func pulse(w http.ResponseWriter, r *http.Request) {
	customerID, videoID := parseQuery(r.URL.RawQuery)

	sessionTimer := time.NewTimer(6 * time.Second)

	go func() {
		<-sessionTimer.C
		deleteSession(customerID, videoID)
	}()

	mutex := &sync.Mutex{}
	mutex.Lock()
	storeSession(customerID, videoID)
	mutex.Unlock()

	log.Printf("pulse with customer %d and video %d", customerID, videoID)
}

func storeSession(customerID, videoID int) {
	if customers[customerID] == nil {
		customers[customerID] = make(map[int]int)
	}

	customers[customerID][videoID]++

	if videos[videoID] == nil {
		videos[videoID] = make(map[int]int)
	}

	videos[videoID][customerID]++
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
	count := len(customers[customerID])
	mutex.Unlock()

	log.Println("customer stat called: ", customerID, count)
}

func videoCount(w http.ResponseWriter, r *http.Request) {
	mutex := &sync.Mutex{}
	videoID := parseIDFromURL(r.URL.Path)

	mutex.Lock()
	count := len(videos[videoID])
	mutex.Unlock()

	log.Println("video stat called: ", videoID, count)
}

func parseIDFromURL(path string) (id int) {
	stringSlice := strings.Split(path, "/")
	id, _ = strconv.Atoi(stringSlice[2])
	return
}

func deleteSession(customerID, videoID int) {
	mutex := &sync.Mutex{}

	mutex.Lock()

	customers[customerID][videoID]--
	if customers[customerID][videoID] < 1 {
		delete(customers[customerID], videoID)
	}
	if len(customers[customerID]) == 0 {
		delete(customers, customerID)
	}

	videos[videoID][customerID]--
	if videos[videoID][customerID] < 1 {
		delete(videos[videoID], customerID)
	}
	if len(videos[videoID]) == 0 {
		delete(videos, videoID)
	}

	mutex.Unlock()
	log.Println(customers)
	log.Println(videos)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/pulse", pulse)
	mux.HandleFunc("/videos/", videoCount)
	mux.HandleFunc("/customers/", customerCount)

	http.ListenAndServe(":9292", mux)
}
