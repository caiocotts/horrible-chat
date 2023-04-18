package main

import (
	"embed"
	"log"
)

//go:embed assets
var fsys embed.FS

func Asset(path string) []byte {
	data, err := fsys.ReadFile("assets/" + path)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
