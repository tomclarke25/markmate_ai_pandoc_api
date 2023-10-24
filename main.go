package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	BearerPrefix    = "Bearer "
	UnauthorizedMsg = "Unauthorized"
	ConversionFail  = "Conversion failed"
)

type RequestBody struct {
	Markdown string `json:"markdown"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, BearerPrefix)

	if token == "" {
		http.Error(w, UnauthorizedMsg, http.StatusUnauthorized)
		return
	}

	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	docxContent, err := ConvertMarkdownToDocx(body.Markdown)
	if err != nil {
		log.Printf("Failed to convert markdown to docx: %v", err)
		http.Error(w, ConversionFail, http.StatusInternalServerError)
		return
	}

	encodedContent := base64.StdEncoding.EncodeToString(docxContent)

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Write([]byte(encodedContent))
}

func ConvertMarkdownToDocx(markdown string) ([]byte, error) {
	if markdown == "" {
		return nil, errors.New("input markdown string is empty")
	}
	cmd := exec.Command("pandoc", "-f", "markdown", "-t", "docx", "-o", "-")
	cmd.Stdin = strings.NewReader(markdown)
	return cmd.Output()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace with the port you're using
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(r)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}