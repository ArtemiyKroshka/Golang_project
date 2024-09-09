package app

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

func setupServerRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.ViewHanlder)
	mux.HandleFunc("/new", handlers.NewHandler)
	mux.HandleFunc("/create", handlers.CreateHandler)
	mux.HandleFunc("/delete", handlers.DeleteHandler)

	return mux
}

func newServer(port *string) *http.Server {
	server := &http.Server{
		Addr:    ":" + *port,
		Handler: setupServerRoutes(),
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

func endServer(server *http.Server, timeout time.Duration) {
	sigterm := make(chan os.Signal, 1) // package "closer" as an alternative
	signal.Notify(sigterm, os.Interrupt)

	<-sigterm

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}

func Run() {
	// Чтение флага для порта
	portFlag := flag.String("port", "8080", "Server port")
	flag.Parse()

	// Инициализация сервера
	serv := newServer(portFlag)

	// Запуск сервера
	startServer(serv)

	// Завершение работы сервера с таймаутом
	endServer(serv, 5*time.Second)
}
