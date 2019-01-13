package main

import (
	"encoding/json"
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

func pulse(w http.ResponseWriter, r *http.Request) {
	customerID, videoID := parseQuery(r.URL)
	if customerID == 0 || videoID == 0 {
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
		return
	}

	sessionTimer := time.NewTimer(6 * time.Second)

	go func() {
		<-sessionTimer.C
		deleteSession(customerID, videoID)
	}()

	storeSession(customerID, videoID)

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

func customerCount(w http.ResponseWriter, r *http.Request) {
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

func videoCount(w http.ResponseWriter, r *http.Request) {
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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/pulse", pulse)
	mux.HandleFunc("/videos/", videoCount)
	mux.HandleFunc("/customers/", customerCount)

	log.Fatal(http.ListenAndServe(":9292", mux))
}
