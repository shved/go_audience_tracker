package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var urlExamples = []string{
	"/path/1/somewhere?video_id=1&customer_id=2",
	"/path?customer_id=2&video_id=1",
	"/path/1/somewhere?video_id=&customer_id=",
	"/path/1/somewhere?ffff=1&asdf=2",
}

func ExampleParseIDFromURL() {
	fmt.Println(parseIDFromURL("/path/1"))
	fmt.Println(parseIDFromURL("/path/23"))
	fmt.Println(parseIDFromURL("path/321"))
	fmt.Println(parseIDFromURL("path/oops/123"))
	// Output:
	// 1
	// 23
	// 321
	// 123
}

func ExampleParseQuery() {
	for _, str := range urlExamples {
		urlExample, _ := url.Parse(str)
		one, two := parseQuery(urlExample)
		fmt.Println(one, two)
	}
	// Output:
	// 2 1
	// 2 1
	// 0 0
	// 0 0
}

func TestDeleteSession(t *testing.T) {
	state.customers = make(map[int]map[int]int)
	state.customers[1] = make(map[int]int)
	state.customers[1][1]++

	state.videos = make(map[int]map[int]int)
	state.videos[1] = make(map[int]int)
	state.videos[1][1]++

	expectedLen := 0

	deleteSession(1, 1)
	if len(state.customers) != expectedLen {
		t.Fatalf("state customers length is %d, expected to be %d", len(state.customers), expectedLen)
	}
	if len(state.videos) != expectedLen {
		t.Fatalf("state videos length is %d, expected to be %d", len(state.videos), expectedLen)
	}
}

func TestStoreSession(t *testing.T) {
	state.customers = make(map[int]map[int]int)
	state.videos = make(map[int]map[int]int)

	storeSession(1, 1)
	storeSession(1, 2)

	if len(state.customers) != 1 {
		t.Fatalf("state customers length is %d, expected to be %d", len(state.customers), 1)
	}
	if len(state.videos) != 2 {
		t.Fatalf("state videos length is %d, expected to be %d", len(state.videos), 2)
	}
}

func TestCustomerCountHandler(t *testing.T) {
	state.customers = make(map[int]map[int]int)
	state.customers[1] = make(map[int]int)
	state.customers[1][1]++
	state.customers[2] = make(map[int]int)
	state.customers[2][1]++
	state.customers[2][2]++

	state.videos = make(map[int]map[int]int)
	state.videos[1] = make(map[int]int)
	state.videos[1][1]++
	state.videos[1][2]++

	req, err := http.NewRequest("GET", "/customers/2", nil)
	checkError(err, t)

	rr := httptest.NewRecorder()

	http.HandlerFunc(customerCountHandler).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	expected := string(`{"count":2}`)
	if rr.Body.String() != expected {
		t.Errorf("Response body differs. Expected %s .\n Got %s instead", expected, rr.Body.String())
	}
}

func TestVideoCountHandler(t *testing.T) {
	state.customers = make(map[int]map[int]int)
	state.customers[1] = make(map[int]int)
	state.customers[1][1]++
	state.customers[2] = make(map[int]int)
	state.customers[2][1]++

	state.videos = make(map[int]map[int]int)
	state.videos[1] = make(map[int]int)
	state.videos[1][1]++
	state.videos[1][2]++

	req, err := http.NewRequest("GET", "/videos/1", nil)
	checkError(err, t)

	rr := httptest.NewRecorder()

	http.HandlerFunc(videoCountHandler).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	expected := string(`{"count":2}`)
	if rr.Body.String() != expected {
		t.Errorf("Response body differs. Expected %s .\n Got %s instead", expected, rr.Body.String())
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
}
