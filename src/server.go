package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func serveIndex(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "../index.html")
}

func handleReq(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	// kinda hacky way to get .json
	http.ServeFile(w, req, "../data/"+req.URL.Path+".json")

	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("REQ: %s %s %s", req.Method, req.URL, elapsed)
}

func main() {
	mux := http.NewServeMux()

	// port number supplyed from command args
	port := ":" + os.Args[1]
	fmt.Println("Go JSON server listening at localhost" + port)

	mux.HandleFunc("/", serveIndex)

	// get all json endpoints from 'data'
	endpoints, err := ioutil.ReadDir("../data")
	if err != nil {
		fmt.Println("Error reading data ...")
		log.Fatal(err)
	}

	fmt.Println("Available Routes:")
	for _, file := range endpoints {
		filename := file.Name()
		extension := filepath.Ext(filename)
		endpoint := filename[0 : len(filename)-len(extension)] // filename without extension

		fmt.Printf("/%s \n", endpoint)
		mux.HandleFunc("/"+endpoint, handleReq)
	}

	log.Fatal(http.ListenAndServe(port, mux))
}
