package main

import (
	"github.com/mrlauer/meteor"
	"net/http"
)

func main() {
	server := meteor.NewServer()
	server.HandleHTTP("/sockjs")
	http.ListenAndServe(":5678", nil)
}
