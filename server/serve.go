package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	fs := http.FileServer(http.Dir(filepath.Join(filepath.Dir("."), "public")))
	http.Handle("/", fs)

	fmt.Println("Serving at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
