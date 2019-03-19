package main

import (
	"fmt"
	"log"
	"net/http"
)

func httpListener() {
	http.HandleFunc("/", httpStatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func httpStatHandler(w http.ResponseWriter, r *http.Request) {
	stat := r.URL.Path[1:]
	fmt.Fprintf(w, `{"count": %d}`, cache.get(stat))
}
