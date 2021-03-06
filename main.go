/*
Audience tracker is a service for tracking client devcies via heartbeats (or pulses).
It registers every get request comes into pulse path with customer_id and video_id parameter
and stores it in state struct for six seconds.
State has two separate maps of maps for customers and videos which are storing session counters for each
customer/video pairs.
Videos/customers endpoints then counts the needed map keys and returning corresponding counters.
*/
package main

import (
	"encoding/json"
	_ "github.com/y0ssar1an/q"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var state = struct {
	sync.Mutex
	customers map[int]map[int]int
	videos    map[int]map[int]int
}{
	customers: make(map[int]map[int]int),
	videos:    make(map[int]map[int]int),
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/pulse", pulseHandler)
	mux.HandleFunc("/videos/", videoCountHandler)
	mux.HandleFunc("/customers/", customerCountHandler)

	log.Fatal(http.ListenAndServe(":9292", mux))
}

func pulseHandler(w http.ResponseWriter, r *http.Request) {
	customerID, videoID := parseQuery(r.URL)
	if customerID == 0 || videoID == 0 {
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
		return
	}

	sessionTimer := time.NewTimer(6 * time.Second)

	go func(customerID, videoID int) {
		storeSession(customerID, videoID)
		<-sessionTimer.C
		deleteSession(customerID, videoID)
	}(customerID, videoID)

	w.WriteHeader(200)

	log.Println("pulse registered for:", customerID, videoID)
}

func storeSession(customerID, videoID int) {
	state.Lock()

	if state.customers[customerID] == nil {
		state.customers[customerID] = make(map[int]int)
	}
	state.customers[customerID][videoID]++
	if state.videos[videoID] == nil {
		state.videos[videoID] = make(map[int]int)
	}
	state.videos[videoID][customerID]++

	state.Unlock()
}

func parseQuery(urlObj *url.URL) (customerID, videoID int) {
	values := urlObj.Query()
	customerID, _ = strconv.Atoi(values.Get("customer_id"))
	videoID, _ = strconv.Atoi(values.Get("video_id"))
	return
}

func customerCountHandler(w http.ResponseWriter, r *http.Request) {
	customerID := parseIDFromURL(r.URL.Path)
	if customerID == 0 {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
		return
	}

	state.Lock()
	count := len(state.customers[customerID])
	state.Unlock()

	responseBody, err := json.Marshal(map[string]interface{}{
		"count": count,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)

	log.Printf("customer %d videos count %d", customerID, count)
}

func videoCountHandler(w http.ResponseWriter, r *http.Request) {
	videoID := parseIDFromURL(r.URL.Path)
	if videoID == 0 {
		http.Error(w, "invalid video id", http.StatusBadRequest)
		return
	}

	state.Lock()
	count := len(state.videos[videoID])
	state.Unlock()

	responseBody, err := json.Marshal(map[string]interface{}{
		"count": count,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)

	log.Printf("video %d customers count %d", videoID, count)
}

func parseIDFromURL(path string) (id int) {
	stringSlice := strings.Split(path, "/")
	id, _ = strconv.Atoi(stringSlice[len(stringSlice)-1])
	return
}

func deleteSession(customerID, videoID int) {
	state.Lock()

	state.customers[customerID][videoID]--
	if state.customers[customerID][videoID] < 1 {
		delete(state.customers[customerID], videoID)
	}
	if len(state.customers[customerID]) == 0 {
		delete(state.customers, customerID)
	}
	state.videos[videoID][customerID]--
	if state.videos[videoID][customerID] < 1 {
		delete(state.videos[videoID], customerID)
	}
	if len(state.videos[videoID]) == 0 {
		delete(state.videos, videoID)
	}

	state.Unlock()

	log.Println("session deleted", customerID, videoID)
}
