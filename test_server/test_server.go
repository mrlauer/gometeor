package main

import (
	"flag"
	"fmt"
	"github.com/mrlauer/meteor"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func startServer() {
	server := meteor.NewServer()
	server.Methods(map[string]interface{}{
		"Greeting": func(arg string) string {
			return fmt.Sprintf("Hello from Go, %s!", arg)
		},
	})
	server.HandleHTTP("/sockjs")
	addr := "127.0.0.1:3010"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("Go server listening on %s\n", addr)
	log.Print(http.Serve(l, nil))
}

func startMeteor() {
	// Get the directory
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get file")
	}
	dir := path.Dir(path.Dir(file))
	meteorDir := path.Join(dir, "test_meteor_project")
	err := os.Chdir(meteorDir)
	if err != nil {
		log.Print(err)
		return
	}
	cmd := exec.Command("meteor")
	fmt.Printf("Starting meteor server on port 3000\n")
	err = cmd.Run()
	if err != nil {
		log.Printf("Error running meteor: %v\n", err)
	}
}

func main() {
	serverOnly := flag.Bool("nometeor", false, "Just run the go server, not the meteor app")
	flag.Parse()
	if *serverOnly {
		startServer()
	} else {
		go startServer()
		startMeteor()
	}
}
