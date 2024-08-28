package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var mu sync.Mutex

type Guestbook struct {
	SignatureCount int
	Signatures     []string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func serverError(err error, writer http.ResponseWriter) {
	if err != nil {
		log.Printf("Error text: %v", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	check(err)

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	check(scanner.Err())

	return lines

}

func viewHanlder(writer http.ResponseWriter, _ *http.Request) {
	signatures := getStrings("src/signatures.txt")
	html, err := template.ParseFiles("src/templates/view.html")
	serverError(err, writer)
	err = html.Execute(writer, Guestbook{Signatures: signatures, SignatureCount: len(signatures)})
	serverError(err, writer)
}

func newHandler(writer http.ResponseWriter, _ *http.Request) {
	html, err := template.ParseFiles("src/templates/new.html")
	serverError(err, writer)
	err = html.Execute(writer, nil)
	serverError(err, writer)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	signature := request.FormValue("signature")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE

	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile("src/signatures.txt", options, os.FileMode(0600))
	serverError(err, writer)
	defer file.Close()
	_, err = fmt.Fprintln(file, signature)
	serverError(err, writer)
	http.Redirect(writer, request, "/", http.StatusFound)
}

func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	index, err := strconv.Atoi(request.FormValue("index"))
	serverError(err, writer)
	signatures := getStrings("src/signatures.txt")

	mu.Lock()
	defer mu.Unlock()

	newSignatures := make([]string, 0, len(signatures)-1)
	options := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	for i, v := range signatures {
		if i == index {
			continue
		}
		newSignatures = append(newSignatures, v)
	}
	file, err := os.OpenFile("src/signatures.txt", options, os.FileMode(0600))
	serverError(err, writer)
	defer file.Close()
	for _, line := range newSignatures {
		_, err := file.WriteString(line + "\n")
		check(err)
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}

func handlerError() {
	p := recover()
	if p == nil {
		return
	}
	err, ok := p.(error)
	if ok {
		fmt.Println("Recovered from panic:", err)
	} else {
		fmt.Println("Recovered from panic:", p)
	}
}

func main() {

	defer handlerError()

	var port = os.Getenv("PORT")

	if port == "" {
		portFlag := flag.String("port", "8080", "Server port")
		flag.Parse()
		port = *portFlag
	}

	http.HandleFunc("/", viewHanlder)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(":"+port, nil)
}
