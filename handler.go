package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Route URIs
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{"root", "GET", "/", EndpointRoot},
	Route{"submit", "GET", "/v1/queue.main/submit", TodoIndex},
	Route{"SubmitJob", "POST", "/v1/queue.main/new", TodoCreate},
	Route{"TodoShow", "GET", "/v1/queue.main/{todoId}", TodoShow},
}

// EndpointRoot returns status and info
func EndpointRoot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&struct {
		Application     string   `json:"application"`
		Version         string   `json:"version"`
		Ready           bool     `json:"ready"`
		EndpointVersion int      `json:"endpoint_version"`
		Endpoint        []string `json:"endpoint"`
	}{
		"Flock - Message queue",
		"1.0.0",
		true,
		1,
		[]string{
			"/v1/queue/list",
			"/v1/queue/create",
			"/v1/queue/remove",
			"/v1/queue.{queue}/submit",
			"/v1/queue.{queue}/purge",
			"/v1/queue.{queue}/{uuid}",
			"/v1/queue.{queue}/status",
		},
	})
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	job := Job{
		UUID:      id,
		Name:      "Some job",
		Completed: false,
		Priority:  12,
		Status:    "SUBMITTED",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(job); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["todoID"]
	fmt.Fprintln(w, "Todo show:", todoID)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(`{"Error": true}`)
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var job Job

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &job); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// t := RepoCreateTodo(job)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

// RequestHandler logs requires and appends required headers
func RequestHandler(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Set server header
		w.Header().Add("Server", "Flock/0.1")

		// Serve reponse
		inner.ServeHTTP(w, r)

		log.Printf("%s -> %s\t%s\t%s\t%s",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

// RESTService starts job listener
func RESTService() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler = RequestHandler(route.HandlerFunc, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
