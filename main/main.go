package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gophercises/urlshort/db"
	"github.com/gophercises/urlshort/urlshort"
)

const (
	boltFile   string = "./bolt.db"
	boltBucket string = "url"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
    - path: /urlshort
      url: https://github.com/gophercises/urlshort
    - path: /urlshort-final
      url: https://github.com/gophercises/urlshort/tree/solution
    `
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		fmt.Printf("error when parsing the yaml file - %s", err)
		os.Exit(1)
	}

	// Create BoltDB
	boltDB, err := db.NewBoltDB(boltFile, 0600, boltBucket)
	if err != nil {
		fmt.Printf("error when creating/reading the bold db - %s", err)
		os.Exit(1)
	}

	// Add some sample records
	boltDB.Insert("/google", "https://www.google.es")
	boltDB.Insert("/amazon", "https://www.amazon.es")

	// Create DBHandler
	dbHandler := urlshort.DBHandler(boltDB, yamlHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
