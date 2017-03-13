package main

import (
	"encoding/json"
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

type ServiceError struct {
	Message   string `json:"message"`
	Solution  string `json:"solution"`
	ErrorCode int    `json:"error_code"`
}

var routes = []Route{
	Route{"root", "GET", "/", EndpointRoot},

	// Endpoints for API version 1
	Route{"v1.root", "GET", "/v1/", v1EndpointList},
	Route{"v1.job.status", "GET", "/v1/queue.{queue}/status", v1JobStatus},
	Route{"SubmitJob", "POST", "/v1/queue.{queue}/submit", v1JobSubmit},
	// Route{"TodoShow", "GET", "/v1/queue.main/{todoId}", TodoShow},
}

func v1JobStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	//TODO: Mock return
	job := &Job{
		UUID:      id,
		Name:      "Some job",
		Version:   1,
		Priority:  10,
		Completed: false,
		Status:    "SUBMITTED",
		Owner:     "MAIN",
		Submitted: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(job); err != nil {
		panic(err)
	}
}

// func TodoShow(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	todoID := vars["todoID"]
// 	fmt.Fprintln(w, "Todo show:", todoID)
// }

func v1JobSubmit(w http.ResponseWriter, r *http.Request) {
	var job Job

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &job); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// t := RepoCreateTodo(job)
	// StoreNewJob()

	// log.Print("heiro")
	// EndpointError(w, &ServiceError{
	// 	Message:   "Something false",
	// 	Solution:  "Make it true",
	// 	ErrorCode: http.StatusBadRequest,
	// })
	// return

	w.WriteHeader(http.StatusCreated)
}

// EndpointError returns an error structure
func EndpointError(w http.ResponseWriter, e *ServiceError) {
	w.WriteHeader(e.ErrorCode)
	json.NewEncoder(w).Encode(e)
}

// EndpointRoot returns status and info
func EndpointRoot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&struct {
		Application         string `json:"application"`
		Version             string `json:"version"`
		Ready               bool   `json:"ready"`
		LastEndpointVersion int    `json:"last_endpoint_version"`
		LastEndpointPrefix  string `json:"last_endpoint_prefix"`
	}{
		"Flock - Message queue",
		version,
		true,
		1,
		"/v1/",
	})
}

// v1EndpointList show enpoints for this version
func v1EndpointList(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&struct {
		EndpointVersion int      `json:"endpoint_version"`
		Endpoint        []string `json:"endpoint"`
	}{
		1,
		[]string{
			"/v1/queue/list",
			"/v1/queue/create",
			"/v1/queue/remove",
			"/v1/queue.{queue}/submit",
			"/v1/queue.{queue}/purge",
			"/v1/queue.{queue}/fetch",
			"/v1/queue.{queue}/status",
			"/v1/job.{uuid}/status",
			"/v1/job.{uuid}/cancel",
			"/v1/task.{uuid}/status",
			"/v1/task.{uuid}/cancel",
		},
	})
}

// NotFound returns notfound erndpoint
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&ServiceError{
		Message:   "Endpoint not found",
		Solution:  "See / for possible directives",
		ErrorCode: http.StatusNotFound,
	})
}

// RequestHandler logs requires and appends required headers
func RequestHandler(inner http.Handler, name string, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		err := false

		// Set server header
		w.Header().Add("Server", "Flock/"+version)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(&ServiceError{
				Message:   "Endpoint does not accept method",
				Solution:  "Change HTTP method",
				ErrorCode: http.StatusMethodNotAllowed,
			})
			err = true
		}

		// Serve reponse
		if !err {
			inner.ServeHTTP(w, r)
		}

		if verbose {
			log.Printf("%s -> %s\t%s\t%s\t%s",
				r.RemoteAddr,
				r.Method,
				r.RequestURI,
				name,
				time.Since(start),
			)
		}
	})
}

// RESTService starts job listener
func RESTService() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(NotFound)

	for _, route := range routes {
		var handler = RequestHandler(route.HandlerFunc, route.Name, route.Method)

		router.
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
