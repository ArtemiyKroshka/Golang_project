package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Guestbook struct {
	SignatureCount int
	Signatures     []string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
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

func viewHanlder(writer http.ResponseWriter, request *http.Request) {
	signatures := getStrings("src/signatures.txt")
	html, err := template.ParseFiles("src/view.html")
	check(err)
	err = html.Execute(writer, Guestbook{Signatures: signatures, SignatureCount: len(signatures)})
	check(err)
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("src/new.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	signature := request.FormValue("signature")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("src/signatures.txt", options, os.FileMode(0600))
	check(err)
	_, err = fmt.Fprintln(file, signature)
	check(err)
	err = file.Close()
	check(err)
	http.Redirect(writer, request, "/guestbook", http.StatusFound)
}

func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	index, err := strconv.Atoi(request.FormValue("index"))
	check(err)
	signatures := getStrings("src/signatures.txt")
	newSignatures := make([]string, 0, len(signatures)-1)
	options := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	for i, v := range signatures {
		if i == index {
			continue
		}
		newSignatures = append(newSignatures, v)
	}
	file, err := os.OpenFile("src/signatures.txt", options, os.FileMode(0600))
	check(err)
	defer file.Close()
	for _, line := range newSignatures {
		_, err := file.WriteString(line + "\n")
		check(err)
	}
	http.Redirect(writer, request, "/guestbook", http.StatusFound)
}

func main() {

	http.HandleFunc("/guestbook", viewHanlder)
	http.HandleFunc("/guestbook/new", newHandler)
	http.HandleFunc("/guestbook/create", createHandler)
	http.HandleFunc("/guestbook/delete", deleteHandler)
	err := http.ListenAndServe("localhost:3000", nil)
	log.Fatal(err)
}
