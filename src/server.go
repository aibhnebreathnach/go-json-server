package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

func serveIndex(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "../index.html")
}

func handleGET(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	// kinda hacky way to get .json
	http.ServeFile(w, req, "../data/"+req.URL.Path+".json")

	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("REQ: %s %s %s", req.Method, req.URL, elapsed)
}

func handlePUT(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	type ReqBody struct {
		JSON string
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var j ReqBody
	err = json.Unmarshal(body, &j)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	elapsed := now.Sub(start)
	log.Printf("REQ: %s %s %s", req.Method, req.URL, elapsed)
}

func handleRequest(w http.ResponseWriter, req *http.Request) {

	switch req.Method {

	case "GET":
		handleGET(w, req)
	case "PUT":
		handlePUT(w, req)
	}
}

func main() {
	r := mux.NewRouter()

	// port number supplyed as commandline arg 1
	port := ":" + os.Args[1]
	fmt.Println("Go JSON server listening at localhost" + port)

	r.HandleFunc("/", serveIndex)

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
		endpoint := "/" + filename[0:len(filename)-len(extension)] // filename without extension

		fmt.Printf("%s \n", endpoint)
		r.HandleFunc(endpoint, handleRequest).Methods("GET", "PUT")

	}

	http.ListenAndServe(port, r)
}
