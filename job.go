package main

import "time"

// Job model
type Job struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Priority  int8      `json:"priority"`
	Status    string    `json:"status"`
	Completed bool      `json:"completed"`
	Submitted time.Time `json:"submitted"`
}

// Jobs array
type Jobs []Job
