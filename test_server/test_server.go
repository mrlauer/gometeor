package main

import (
	"fmt"
	"github.com/mrlauer/meteor"
	"net/http"
)

func main() {
	server := meteor.NewServer()
	server.RegisterFunction("FooMethod", func(arg string) string {
		return fmt.Sprintf("Hello from Go, %s!", arg)
	})
	server.HandleHTTP("/sockjs")
	http.ListenAndServe(":5678", nil)
}
