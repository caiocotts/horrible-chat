package main

import (
	_ "embed"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"log"
	"net/http"
	"path"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("New connection")
		w.Write(Asset("index.html"))
	})

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		id, err := gonanoid.New()
		if err != nil {
			log.Println(err)
			return
		}
		http.Redirect(w, r, "/chat/"+id, http.StatusFound)
	})

	http.HandleFunc("/chat/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(Asset("chat.html"))
		_, id := path.Split(r.URL.Path)
		log.Printf("Your session is %s\n", id)
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:4533", nil))
}
