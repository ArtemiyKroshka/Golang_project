package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

var mu sync.RWMutex

type Guestbook struct {
	SignatureCount int
	Signatures     []string
}

func serverError(writer http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("Error text: %v", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getStrings(fileName string) []string {

	mu.Lock()
	defer mu.Unlock()

	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		log.Println("Error opening file:", err)
		return nil
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
	}

	return lines

}

func viewHanlder(writer http.ResponseWriter, _ *http.Request) {
	signatures := getStrings("src/signatures.txt")
	html, err := template.ParseFiles("src/templates/view.html")
	if err != nil {
		serverError(writer, err)
		return
	}
	if err := html.Execute(writer, Guestbook{Signatures: signatures, SignatureCount: len(signatures)}); err != nil {
		serverError(writer, err)
	}
}

func newHandler(writer http.ResponseWriter, _ *http.Request) {
	html, err := template.ParseFiles("src/templates/new.html")
	if err != nil {
		serverError(writer, err)
		return
	}
	if err := html.Execute(writer, nil); err != nil {
		serverError(writer, err)
	}
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	signature := request.FormValue("signature")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE

	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile("src/signatures.txt", options, os.FileMode(0600))
	if err != nil {
		serverError(writer, err)
		return
	}
	defer file.Close()
	if _, err := fmt.Fprintln(file, signature); err != nil {
		serverError(writer, err)
		return
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}

func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	index, err := strconv.Atoi(request.FormValue("index"))
	if err != nil {
		serverError(writer, err)
		return
	}
	signatures := getStrings("src/signatures.txt")

	mu.Lock()
	defer mu.Unlock()

	if index < 0 || index >= len(signatures) {
		http.Error(writer, "Invalid index", http.StatusBadRequest)
		return
	}

	newSignatures := make([]string, 0, len(signatures)-1)
	options := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	for i, v := range signatures {
		if i == index {
			continue
		}
		newSignatures = append(newSignatures, v)
	}
	file, err := os.OpenFile("src/signatures.txt", options, os.FileMode(0600))
	if err != nil {
		serverError(writer, err)
		return
	}
	defer file.Close()
	for _, line := range newSignatures {
		if _, err := file.WriteString(line + "\n"); err != nil {
			serverError(writer, err)
			return
		}
	}
	http.Redirect(writer, request, "/", http.StatusFound)
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
		Addr: ":" + *portFlag,
	}

	http.HandleFunc("/", viewHanlder)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/delete", deleteHandler)

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
