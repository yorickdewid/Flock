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

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Queue struct {
	Name       string `json:"name"`
	JobsCount  int    `json:"jobs_count"`
	JobsDone   int    `json:"jobs_done"`
	JobsAwait  int    `json:"jobs_await"`
	JobsFailed int    `json:"jobs_failed"`
}

func SetupStore() {
	// Bail if store exists
	if _, err := os.Stat("job.store"); err == nil {
		return
	}

	db, err := sql.Open("sqlite3", "job.store")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Print("Setup datastore")
	_, err = db.Exec(`
	CREATE TABLE queue (name TEXT PRIMARY KEY NOT NULL,
						can_remove INT DEFAULT 1,
						jobs_count INT DEFAULT 0,
						jobs_done INT DEFAULT 0,
						jobs_await INT DEFAULT 0,
						jobs_failed INT DEFAULT 0);
	CREATE TABLE job (id TEXT PRIMARY KEY NOT NULL,
						queue TEXT NOT NULL,
						name TEXT,
						version TEXT,
						status TEXT,
						owner TEXT,
						priority INT DEFAULT 10,
						completed INT DEFAULT 0,
						content_file TEXT,
						submitted_at DATETIME,
						updated_at TIMESTAMP);
	CREATE TABLE task (id TEXT PRIMARY KEY NOT NULL,
						job TEXT NOT NULL,
						command TEXT,
						name TEXT,
						arguments TEXT,
						completed INT DEFAULT 0);
	INSERT INTO queue ('name','can_remove') VALUES ('main',0);
	`)
	if err != nil {
		log.Printf("Create database error: %q\n", err)
		return
	}
}

func StoreNewJob() {
	db, err := sql.Open("sqlite3", "job.store")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO job (id,queue,name,submitted_at) VALUES (?, ?, ?, DATETIME('now'))")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var name string
	if _, err := stmt.Exec("f6146500-f0ce-49a5-af1a-c0525ea0ced7", "main", "SomeJob"); err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}

func StoreQueueList() *[]Queue {
	db, err := sql.Open("sqlite3", "job.store")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT name,jobs_count,jobs_done,jobs_await,jobs_failed FROM queue")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var queueList []Queue
	for rows.Next() {
		var queue Queue
		err = rows.Scan(&queue.Name, &queue.JobsCount, &queue.JobsDone, &queue.JobsAwait, &queue.JobsFailed)
		if err != nil {
			log.Fatal(err)
		}

		queueList = append(queueList, queue)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return &queueList
}

func StoreQueueCreate(queue *Queue) error {
	db, err := sql.Open("sqlite3", "job.store")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check empty input
	queue.Name = strings.ToLower(strings.TrimSpace(queue.Name))
	if len(queue.Name) == 0 {
		return &ServiceError{
			Message:  "Name is empty",
			Solution: "Supply a name",
		}
	}

	// Check existing name
	row := db.QueryRow("SELECT 1 FROM queue WHERE name=? LIMIT 1", queue.Name)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	row.Scan(&count)
	if count > 0 {
		return &ServiceError{
			Message:  "Queue with name already exists",
			Solution: "Rename input object",
		}
	}

	stmt, err := db.Prepare("INSERT INTO queue (name) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(queue.Name); err != nil {
		log.Fatal(err)
	}

	log.Print("New queue: ", queue.Name)

	return err
}

func StoreQueueRemove(queue *Queue) error {
	db, err := sql.Open("sqlite3", "job.store")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check empty input
	queue.Name = strings.ToLower(strings.TrimSpace(queue.Name))
	if len(queue.Name) == 0 {
		return &ServiceError{
			Message:  "Name is empty",
			Solution: "Supply a name",
		}
	}

	// Check existing name
	row := db.QueryRow("SELECT 1 FROM queue WHERE name=? AND can_remove=1 LIMIT 1", queue.Name)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	row.Scan(&count)
	if count == 0 {
		return &ServiceError{
			Message:  "Queue with name does not exist",
			Solution: "Rename input object",
		}
	}

	stmt, err := db.Prepare("DELETE FROM queue WHERE name=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(queue.Name); err != nil {
		log.Fatal(err)
	}

	log.Print("remove queue: ", queue.Name)

	return err
}

/*
	////
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
	////

	///
	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	///

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	_, err = db.Exec("delete from foo")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}
*/
