package main

import (
	"go_project/internal/server"
	"log"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	server.Run()

}
