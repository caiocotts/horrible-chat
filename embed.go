package main

import (
	"embed"
	"log"
)

//go:embed assets
var Fsys embed.FS

func init() {
	file, err := Fsys.Open("assets/font/3270.ttf")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func Asset(path string) []byte {
	data, err := Fsys.ReadFile("assets/" + path)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func Font() []byte {
	return Asset("font/3270.ttf")
}
