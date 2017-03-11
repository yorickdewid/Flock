package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var routes = Routes{
	Route{"Index", "GET", "/", Index},
	Route{"TodoIndex", "GET", "/todos", TodoIndex},
	Route{"SubmitJob", "POST", "/todo/new", TodoCreate},
	Route{"TodoShow", "GET", "/todos/{todoId}", TodoShow},
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "{\"error\":\"no jobs\"}")
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
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
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
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// t := RepoCreateTodo(job)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	// if err := json.NewEncoder(w).Encode(t); err != nil {
	// panic(err)
	// }
}
