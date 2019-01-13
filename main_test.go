package main

import (
	"fmt"
	"net/url"
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

// func TestStoreSession(t *testing.T) {
// 	customers[1] = make(map[int]int)
// 	customers[1][1]++
// 	videos[1] = make(map[int]int)
// 	videos[1][1]++
//
// 	deleteSession(1, 1)
// 	if len(sessions) > 0 || len(videos) > 0 {
// 	}
// }
