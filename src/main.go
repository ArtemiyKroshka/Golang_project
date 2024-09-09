package main

import (
	"go_project/src/app"
	"log"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	app.Run()

}
