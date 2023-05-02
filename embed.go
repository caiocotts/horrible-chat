package main

import (
	"embed"
	"log"
)

//go:embed assets
var Fsys embed.FS

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
