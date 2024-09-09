package handlers

import (
	"bytes"
	"fmt"
	"go_project/src/utils/readFile"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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

func ViewHanlder(writer http.ResponseWriter, _ *http.Request) {
	signatures := readFile.GetStrings("src/signatures.txt", &mu)
	html, err := template.ParseFiles("src/templates/view.html")
	if err != nil {
		serverError(writer, err)
		return
	}
	var buf bytes.Buffer
	if err := html.Execute(&buf, Guestbook{Signatures: signatures, SignatureCount: len(signatures)}); err != nil {
		serverError(writer, err)
	}

	if _, err := buf.WriteTo(writer); err != nil {
		log.Printf("Cannot write response for view handler: %v", err)
	}
}

func NewHandler(writer http.ResponseWriter, _ *http.Request) {
	html, err := template.ParseFiles("src/templates/new.html")
	if err != nil {
		serverError(writer, err)
		return
	}
	var buf bytes.Buffer
	if err := html.Execute(&buf, nil); err != nil {
		serverError(writer, err)
	}
	if _, err := buf.WriteTo(writer); err != nil {
		log.Printf("Cannot write response for new handler: %v", err)
	}
}

func CreateHandler(writer http.ResponseWriter, request *http.Request) {
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

func DeleteHandler(writer http.ResponseWriter, request *http.Request) {
	index, err := strconv.Atoi(request.FormValue("index"))
	if err != nil {
		serverError(writer, err)
		return
	}
	signatures := readFile.GetStrings("src/signatures.txt", &mu)

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
