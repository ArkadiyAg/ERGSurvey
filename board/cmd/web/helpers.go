package main

import (
	"embed"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
)

//go:embed templates
var templateFS embed.FS

func renderQuestions(w http.ResponseWriter, t string, questions []string) {
	tmpl, err := template.ParseFS(templateFS, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, questions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// readJSON tries to read the body of a request and converts it into JSON
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
