package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type session struct {
	videoID    int
	customerID int
}

func pulse(w http.ResponseWriter, r *http.Request) {
	log.Printf("pulse registered on %s", r.URL.Path)
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}
	customerID, err := strconv.Atoi(values.Get("customer_id"))
	if err != nil {
		log.Fatal(err)
	}
	videoID, err := strconv.Atoi(values.Get("video_id"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("with customer id %d", customerID)
	log.Printf("with video id %d", videoID)
}

func videoStat(w http.ResponseWriter, r *http.Request) {
	log.Printf("video stat called")
}

func customerStat(w http.ResponseWriter, r *http.Request) {
	log.Printf("customer stat called")
}

func main() {
	// var sessions []session
	mux := http.NewServeMux()

	mux.HandleFunc("/pulse", pulse)
	mux.HandleFunc("/videos/", videoStat)
	mux.HandleFunc("/customers/", customerStat)

	http.ListenAndServe(":9292", mux)
}
