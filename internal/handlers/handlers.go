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
	Signatures     []database.Line
}

func serverError(writer http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("Error text: %v", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ViewHandler(database *database.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		signatures, err := database.GetData()
		if err != nil {
			serverError(writer, err)
			return
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
}

func NewHandler(database *database.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
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
}

func CreateHandler(database *database.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		signature := request.FormValue("signature")
		if err := database.SetData(signature); err != nil {
			serverError(writer, err)
			return
		}

		http.Redirect(writer, request, "/", http.StatusFound)
	}
}

func DeleteHandler(database *database.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := request.FormValue("id")
		if err := database.DeleteData(id); err != nil {
			serverError(writer, err)
			return
		}

		http.Redirect(writer, request, "/", http.StatusFound)
	}
}
