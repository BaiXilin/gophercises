package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/xilin/cyoa"
)

func main() {
	// Creates a command that takes in the name of the story json file
	port := flag.Int("port", 3030, "the port to start CYOA web application on")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	// open the file
	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	// decode the json file into Go struct Story
	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server at %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
