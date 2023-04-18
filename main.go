package main

import (
	_ "embed"
	"log"
	"net/http"
)

var (
	//go:embed index.html
	indexHtml []byte
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		println("New connection")
		w.Write(indexHtml)
	})
	log.Fatal(http.ListenAndServe(":4533", nil))
}
