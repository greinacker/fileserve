package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var rootPath string
var signingToken string

func handler(w http.ResponseWriter, req *http.Request) {

	log.Println(req.Method, req.URL.Path)

	switch req.Method {
	case "GET", "HEAD":
		http.ServeFile(w, req, rootPath+req.URL.Path)
	case "POST", "PUT":
		buf, _ := ioutil.ReadAll(req.Body)

		if signingToken != "" {
			signature := req.Header.Get("FileserveSignature")
			if signature == "" {
				w.WriteHeader(http.StatusForbidden)
				log.Println("Error: unsigned request received from", req.RemoteAddr, req.URL.Path)
				return
			}

			sigParts := strings.Split(signature, ":")
			clientTime, _ := strconv.ParseInt(sigParts[0], 10, 64)

			if (time.Now().Unix()-clientTime) > 180 || (time.Now().Unix()-clientTime) < -180 {
				w.WriteHeader(http.StatusForbidden)
				log.Println("Error: request timestamp too old from", req.RemoteAddr, req.URL.Path)
				return
			}

			h := sha1.New()
			io.WriteString(h, sigParts[0]+signingToken+req.URL.Path)
			h.Write(buf)
			expectSig := fmt.Sprintf("%x", h.Sum(nil))

			if expectSig != sigParts[1] {
				w.WriteHeader(http.StatusForbidden)
				log.Println("Error: invalid signature from", req.RemoteAddr, req.URL.Path)
				return
			}
		}

		err := ioutil.WriteFile(rootPath+req.URL.Path, buf, 0644)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error:", err)
		}

	default:

	}
}

func main() {
	signingToken = strings.TrimRight(os.Getenv("SIGN_SECRET"), "")
	if signingToken == "" {
		log.Println("Warning: SIGN_SECRET not found, all writes will succeed without signature")
	} else {
		log.Println("SIGN_SECRET found, all writes will require signature")
	}

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
