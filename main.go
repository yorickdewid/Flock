package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	build   = "none"
	version = "dev-build"
)

func acceptNewJobs(port *string) {
	setupStore()
	log.Print("Waiting for jobs on :" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, RESTService()))
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
	// var configFlag = flag.String("c", "", "config file")
	flag.Parse()

	if *versionFlag {
		printVersion()
		return
	}

	if *helpFlag {
		printHelp()
		return
	}

	acceptNewJobs(portFlag)
}
