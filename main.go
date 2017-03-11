package main

import (
	"log"
	"net/http"
)

func acceptNewJobs() {
	log.Print("Waiting for jobs on :44087")
	log.Fatal(http.ListenAndServe(":44087", RESTService()))
}

func main() {
	acceptNewJobs()
}
