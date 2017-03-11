package main

import (
	"log"
	"net/http"
)

func listenForJobs() {
	log.Print("Waiting for jobs on :44087")
	log.Fatal(http.ListenAndServe(":44087", RESTService()))
}

func main() {
	listenForJobs()
}
