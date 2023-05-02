package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	broker := newBroker()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(Asset("index.html"))
	})

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		id, err := gonanoid.New()
		if err != nil {
			log.Println(err)
			return
		}
		broker.newChats <- id
		http.Redirect(w, r, "/c/"+id, http.StatusFound)
	})

	http.HandleFunc("/c/", func(w http.ResponseWriter, r *http.Request) {
		_, id := path.Split(r.URL.Path)
		if _, ok := broker.chats[id]; !ok {
			//http.Error(w, "this chat does not exist", http.StatusNotFound)
			_, _ = w.Write(Asset("notfound.html"))
			return
		}
		t, _ := template.New("").Parse(string(Asset("chat.html")))
		uid, _ := gonanoid.New()
		_ = t.Execute(w, struct {
			UserId string
		}{UserId: uid})

		log.Printf("Your session is %s\n", id)
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			var m message

			dec := json.NewDecoder(r.Body)
			dec.DisallowUnknownFields()
			err := dec.Decode(&m)
			if err != nil {

				var syntaxError *json.SyntaxError
				var unmarshalTypeError *json.UnmarshalTypeError

				switch {
				case errors.As(err, &syntaxError):
					msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
					http.Error(w, msg, http.StatusBadRequest)

				case errors.Is(err, io.ErrUnexpectedEOF):
					msg := "Request body contains badly-formed JSON"
					http.Error(w, msg, http.StatusBadRequest)

				case errors.As(err, &unmarshalTypeError):
					msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
					http.Error(w, msg, http.StatusBadRequest)

				case strings.HasPrefix(err.Error(), "json: unknown field "):
					fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
					msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
					http.Error(w, msg, http.StatusBadRequest)

				case errors.Is(err, io.EOF):
					msg := "Request body must not be empty"
					http.Error(w, msg, http.StatusBadRequest)

				case err.Error() == "http: request body too large":
					msg := "Request body must not be larger than 1MB"
					http.Error(w, msg, http.StatusRequestEntityTooLarge)

				default:
					log.Print(err.Error())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				return
			}
			err = dec.Decode(&struct{}{})
			if err != io.EOF {
				msg := "Request body must only contain a single JSON object"
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			if strings.TrimSpace(m.Message) == "" {
				return
			}
			log.Printf("Recieved message \"%s\" from user %s", m.Message, m.UserId)
			broker.incoming <- m
		}
	})

	http.HandleFunc("/font/3270.ttf", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(Font())
		if err != nil {
			log.Println(err)
		}
	})

	http.Handle("/events", broker)
	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "localhost:8080"
	}
	log.Println("Horrible Chat has started on:", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
