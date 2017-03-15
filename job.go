//
// Copyright (C) 2017 Quenza Inc.
// All Rights Reserved
//
// This file is part of the Flock project.
//
// Content can not be copied and/or distributed without the express
// permission of the author.
//

package main

import "time"
import "github.com/google/uuid"

// Task model
type Task struct {
	UUID      uuid.UUID `json:"uuid"`
	Command   string    `json:"command"`
	Arguments []string  `json:"arguments"`
	Completed bool      `json:"completed"`
}

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
	Tasks     []Task    `json:"tasks"`
}
