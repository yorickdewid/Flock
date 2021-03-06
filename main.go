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
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	build   = "none"
	version = "1.0.0-dev"
)

var verbose bool

func acceptNewJobs(port *string) {
	SetupStore()

	if verbose {
		log.Print("Running verbose mode")
	}

	srv := &http.Server{
		Handler:      RESTService(),
		Addr:         ":" + *port,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Print("Waiting for jobs on :" + *port)
	log.Fatal(srv.ListenAndServe())
}

func printHelp() {
	fmt.Println("Flock - Message queue\nUsage of flock:")
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Printf("Flock version %v, build %v\n", version, build)
}

func main() {
	var versionFlag = flag.Bool("v", false, "output version information and exit")
	var helpFlag = flag.Bool("h", false, "display this help dialog")
	var portFlag = flag.String("p", "44087", "daemon port")
	var verboseFlag = flag.Bool("V", true, "verbose output")
	flag.Parse()

	if *versionFlag {
		printVersion()
		return
	}

	if *helpFlag {
		printHelp()
		return
	}

	verbose = *verboseFlag
	acceptNewJobs(portFlag)
}
