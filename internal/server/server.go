package server

import (
	"context"
	"flag"
	database "go_project/internal/db"
	"go_project/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupServerRoutes(database *database.Database) *http.ServeMux {
	mux := http.NewServeMux()

	// Pass the database instance to each handler
	mux.HandleFunc("/", handlers.ViewHandler(database))
	mux.HandleFunc("/new", handlers.NewHandler(database))
	mux.HandleFunc("/create", handlers.CreateHandler(database))
	mux.HandleFunc("/delete", handlers.DeleteHandler(database))

	return mux
}

func newServer(port *string, database *database.Database) *http.Server {
	server := &http.Server{
		Addr:    ":" + *port,
		Handler: setupServerRoutes(database),
	}
	return server
}

func startServer(server *http.Server) {
	go func() {
		log.Printf("Server started on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()
}

func endServer(server *http.Server, database *database.Database, timeout time.Duration) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)

	<-sigterm

	log.Println("Shutdown signal received, exiting...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}
	if err := database.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
	log.Println("Server exited properly")
}

func Run() {
	// Get port value from flag
	portFlag := flag.String("port", "8080", "Server port")
	flag.Parse()

	// Initialize database
	database, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize server with the database
	server := newServer(portFlag, database)

	// Start server
	startServer(server)

	// Graceful shutdown
	endServer(server, database, 5*time.Second)
}
