package main

import (
	"embed"
	"log"
)

//go:embed assets
var Assets embed.FS

//go:embed templates
var Templates embed.FS

func init() {
	file, err := Assets.Open("assets/font/3270.ttf")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func Asset(path string) []byte {
	data, err := Assets.ReadFile("assets/" + path)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func Template(path string) []byte {
	data, err := Templates.ReadFile("templates/" + path)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
