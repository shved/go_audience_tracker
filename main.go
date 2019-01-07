package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type session struct {
	VideoID    uint64
	CustomerID uint64
}

var sessions []session

func pulse(w http.ResponseWriter, r *http.Request) {
	customerID, videoID := parseQuery(r.URL.RawQuery)
	newSession := session{VideoID: videoID, CustomerID: customerID}
	mutex := &sync.Mutex{}

	mutex.Lock()
	sessions = append(sessions, newSession)
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
	log.Println("video stat called: ", len(sessions))
}

func customerStat(w http.ResponseWriter, r *http.Request) {
	log.Println("customer stat called: ", len(sessions))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/pulse", pulse)
	mux.HandleFunc("/videos/", videoStat)
	mux.HandleFunc("/customers/", customerStat)

	http.ListenAndServe(":9292", mux)
}
