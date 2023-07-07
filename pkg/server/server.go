package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	word := extractSearchTerm(r.URL.Path)
	if word == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"word": word,
	})
}

// extractSearchTerm extracts the search term from the URL path. eg. /hello -> hello
func extractSearchTerm(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
