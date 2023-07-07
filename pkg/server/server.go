package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	http.HandleFunc("/search/", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func (s *Server) Stop() {
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	query := extractSearchTerm(r.URL.Path)

	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jishoResponse, err := queryJisho(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jishoResponse))
}

func queryJisho(word string) (string, error) {
	JISHO_URL := "https://jisho.org/api/v1/search/words?keyword=%s"
	url := fmt.Sprintf(JISHO_URL, word)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// if the response is not 200, return an error
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: %s", resp.Status)
	}

	// now we can just convert the response body to a string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
