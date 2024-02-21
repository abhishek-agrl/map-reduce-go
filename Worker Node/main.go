package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"worker-node/workers"
)

func Map(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received to start map function")
	fileName, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Request Body Error: %v", err)
	}
	log.Printf("Map Request Body: %v", string(fileName))
	workers.Mapper(string(fileName))
	fmt.Fprintf(w, "Done Mapping: %v", string(fileName))
}

func Reduce(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received to start reduce function")
	fileName, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Request Body Error: %v", err)
	}
	log.Printf("Reduce Request Body: %v", string(fileName))
	workers.Reducer(strings.Split(strings.Trim(string(fileName), " "), " "))
	fmt.Fprintf(w, "Done Reducing: %v", string(fileName))
}

func main() {
	log.Println("New Worker Node Initiated")
	workers.InitWorker()
	defer workers.CloudStoreClient.Close()

	http.HandleFunc("/map", Map)
	http.HandleFunc("/reduce", Reduce)

	log.Println("Starting New HTTP Server at Port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
