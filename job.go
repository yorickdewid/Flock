package main

import "time"
import "github.com/google/uuid"

// Job model
type Job struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Version   int8      `json:"version"`
	Priority  int8      `json:"priority"`
	Status    string    `json:"status"`
	Completed bool      `json:"completed"`
	Owner     string    `json:"owner"`
	Submitted time.Time `json:"submitted"`
	Updated   time.Time `json:"updated"`
	Content   []byte    `json:"content"`
}
