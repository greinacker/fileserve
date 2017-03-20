package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var rootPath string

func handler(w http.ResponseWriter, req *http.Request) {

	log.Println(req.Method + " " + req.URL.Path)

	switch req.Method {
	case "GET":
		http.ServeFile(w, req, rootPath+req.URL.Path)
	case "POST", "PUT":
		buf, _ := ioutil.ReadAll(req.Body)
		err := ioutil.WriteFile(rootPath+req.URL.Path, buf, 0644)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error:", err)
		}

	default:

	}
}

func main() {
	rootPath = strings.TrimRight(os.Getenv("FILESERVE_ROOT"), "/")
	if rootPath == "" {
		log.Println("FILESERVE_ROOT must be present and point to a subdirectory; exiting.")
	} else {
		log.Printf("Serving path %s", rootPath)
		ipAddr := os.Getenv("FILESERVE_IP")
		port := os.Getenv("FILESERVE_PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Listening on %s:%s", ipAddr, port)
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(ipAddr+":"+port, nil))
	}
}
