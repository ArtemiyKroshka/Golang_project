package handlers

import (
	"bytes"
	"fmt"
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

// Define the handlers as methods of a struct
type Handlers struct {
	Database *database.Database
}

func NewHandlers(database *database.Database) *Handlers {
	return &Handlers{Database: database}
}

func (h *Handlers) ViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	signatures, err := h.Database.GetData(ctx)
	if err != nil {
		serverError(w, err)
		return
	}
	html, err := template.ParseFiles("internal/templates/view.html")
	if err != nil {
		serverError(w, err)
		return
	}
	var buf bytes.Buffer
	if err := html.Execute(&buf, Guestbook{Signatures: signatures, SignatureCount: len(signatures)}); err != nil {
		serverError(w, err)
		return
	}

	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("Cannot write response for view handler: %v", err)
	}
}

func (h *Handlers) NewHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("internal/templates/new.html")
	if err != nil {
		serverError(w, err)
		return
	}
	var buf bytes.Buffer
	if err := html.Execute(&buf, nil); err != nil {
		serverError(w, err)
		return
	}
	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("Cannot write response for new handler: %v", err)
	}
}

func (h *Handlers) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	signature := r.FormValue("signature")
	if err := h.Database.SetData(ctx, signature); err != nil {
		serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.FormValue("id")
	if err := h.Database.DeleteData(ctx, id); err != nil {
		fmt.Println("test")
		serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
