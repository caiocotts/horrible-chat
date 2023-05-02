package main

import (
	"fmt"
	"log"
	"net/http"
)

type message struct {
	ChatId  string `json:"chatId"`
	UserId  string `json:"userId"`
	Message string `json:"message"`
}

type connection struct {
	chatId   string
	userId   string
	outgoing chan message
}

type broker struct {
	incoming           chan message
	newConnections     chan connection
	closingConnections chan connection
	newChats           chan string
	chats              map[string]int
	connections        map[string]connection
}

func newBroker() *broker {
	b := &broker{
		incoming:           make(chan message),
		newConnections:     make(chan connection),
		closingConnections: make(chan connection),
		newChats:           make(chan string),
		chats:              make(map[string]int),
		connections:        make(map[string]connection),
	}
	go b.listen()
	return b
}
func (b *broker) listen() {
	for {
		select {
		case c := <-b.newConnections:
			if _, ok := b.chats[c.chatId]; !ok {
				continue
			}
			b.connections[c.userId] = c
			b.chats[c.chatId] += 1
			log.Printf("user %s has connected\n", c.userId)

		case c := <-b.closingConnections:
			if _, ok := b.chats[c.chatId]; !ok {
				continue
			}
			delete(b.connections, c.userId)
			b.chats[c.chatId] -= 1
			if b.chats[c.chatId] == 0 {
				delete(b.chats, c.chatId)
			}
			log.Printf("user %s has disconnected from chat: %s\n", c.userId, c.chatId)

		case cid := <-b.newChats:
			log.Println("chat", cid, "has been created")
			b.chats[cid] = 0
			log.Println(b.chats)

		case m := <-b.incoming:
			log.Printf("User %s has sent this message \"%s\"", m.UserId, m.Message)
			for c := range b.connections {
				if b.connections[c].chatId == m.ChatId {
					b.connections[c].outgoing <- m
				}
			}
		}
	}
}

func (b *broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var f http.Flusher
	var ok bool
	if f, ok = w.(http.Flusher); !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		fmt.Fprintf(w, "Streaming unsupported!")
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	c := connection{
		chatId:   r.URL.Query().Get("chatId"),
		userId:   r.URL.Query().Get("userId"),
		outgoing: make(chan message),
	}
	b.newConnections <- c

	defer func() { b.closingConnections <- c }()

	for {
		select {
		case m := <-c.outgoing:
			fmt.Fprintf(w, "data: {\"chatId\":\"%s\",\"userId\":\"%s\",\"message\":\"%s\"}\n\n", m.ChatId, m.UserId, m.Message)
			log.Println(m)
			f.Flush()
		case <-r.Context().Done():
			return
		}
	}

}
