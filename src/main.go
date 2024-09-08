package main

import (
	"context"
	"flag"
	"go_project/src/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.ViewHanlder)
	mux.HandleFunc("/new", handlers.NewHandler)
	mux.HandleFunc("/create", handlers.CreateHandler)
	mux.HandleFunc("/delete", handlers.DeleteHandler)

	return mux
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	portFlag := flag.String("port", "8080", "Server port")
	flag.Parse()

	server := &http.Server{
		Addr:    ":" + *portFlag,
		Handler: setupRoutes(),
	}

	go func() {
		log.Printf("Server started on port %s", *portFlag)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Start server error: %v", err)
		}
	}()

	sigterm := make(chan os.Signal, 1) // package "closer" as an alternative
	signal.Notify(sigterm, os.Interrupt)

	<-sigterm

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}
