package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// HTTP server root path handler
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Handle page not found
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	io.WriteString(w, "<h1>Hello, Welcome to my HTTP server!</h1>")
}

// HTTP server success request handler
func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	// Return 200
	w.WriteHeader(200)
	fmt.Printf("Success request at %v\n", time.Now())
}

// HTTP server /headers request handler
func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request at %v\n", time.Now())

	// Get VERSION from env variables and insert into req
	// header
	version := os.Getenv("VERSION")
	r.Header.Add("Version", version)

	// Print request headers to console
	for k, v := range r.Header {
		log.Printf("%v: %v\n", k, v)
	}

	// Return request headers to client
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%v: %v\n", k, v))
	}
}
