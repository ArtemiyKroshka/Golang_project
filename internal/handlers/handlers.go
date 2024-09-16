package handlers

import (
	"bytes"
	database "go_project/internal/db"
	"html/template"
	"log"
	"net/http"
)

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
	signatures, err := database.GetData()
	if err != nil {
		log.Fatal(err)
	}
	html, err := template.ParseFiles("internal/templates/view.html")
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
	html, err := template.ParseFiles("internal/templates/new.html")
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

	database.SetData(signature)

	http.Redirect(writer, request, "/", http.StatusFound)
}

func DeleteHandler(writer http.ResponseWriter, request *http.Request) {
	signature := request.FormValue("signature")

	database.DeleteData(signature)

	http.Redirect(writer, request, "/", http.StatusFound)
}
