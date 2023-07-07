package server

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestExtractSearchTerm(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"/search/", ""},
		{"/search/ ", ""},
		{"/search/  ", ""},
		{"/search/   ", ""},
		{"/search/hello", "hello"},
		{"/search/hello world", "hello world"},
		{"/search/言葉", "言葉"},
		{"/hello", ""},
	}

	for _, test := range tests {
		if output := extractSearchTerm(test.input); output != test.expected {
			t.Errorf("Test Failed: %s inputted, %s expected, received: %s", test.input, test.expected, output)
		}
	}
}

func TestHandleValidRequest(t *testing.T) {
	tests := []string{
		"/search/hello",
		"/search/kotoba",
		"/search/tatoeba",
	}

	var wg sync.WaitGroup
	wg.Add(len(tests))

	for _, test := range tests {
		go func(path string) {
			defer wg.Done()

			handler := http.HandlerFunc(handleRequest)

			// Create a new test HTTP request
			req := httptest.NewRequest("GET", path, nil)

			// Create a response recorder to capture the server's response
			rec := httptest.NewRecorder()

			// Serve the HTTP request and record the response
			handler.ServeHTTP(rec, req)

			// Check the response status code
			if rec.Code != http.StatusOK {
				t.Errorf("path=%s. expected status 200, got %d", path, rec.Code)
			}
		}(test)
	}

	wg.Wait()
}

func TestHandleInvalidRequest(t *testing.T) {
	tests := []string{
		"/%20%20",
		"/search/%20%20",
		"/%20%20%20/%20/",
		"/search/%20%20%20/%20/",
		"/",
		"/search/",
	}

	var wg sync.WaitGroup
	wg.Add(len(tests))

	for _, test := range tests {
		go func(path string) {
			defer wg.Done()

			handler := http.HandlerFunc(handleRequest)

			// Create a new test HTTP request
			req := httptest.NewRequest("GET", path, nil)

			// Create a response recorder to capture the server's response
			rec := httptest.NewRecorder()

			// Serve the HTTP request and record the response
			handler.ServeHTTP(rec, req)

			// Check the response status code
			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", rec.Code)
			}
		}(test)
	}

	wg.Wait()
}

func TestQueryJisho(t *testing.T) {
	tests := []struct {
		Word     string
		Expected string
	}{
		{"tatoeba", "{\"meta\":{\"status\":200},\"data\":[{\"slug\":\"例えば\",\"is_common\":true,\"tags\":[\"wanikani14\"],\"jlpt\":[\"jlpt-n4\"],\"japanese\":[{\"word\":\"例えば\",\"reading\":\"たとえば\"}],\"senses\":[{\"english_definitions\":[\"for example\",\"for instance\",\"e.g.\"],\"parts_of_speech\":[\"Adverb (fukushi)\"],\"links\":[],\"tags\":[],\"restrictions\":[],\"see_also\":[],\"antonyms\":[],\"source\":[],\"info\":[]}],\"attribution\":{\"jmdict\":true,\"jmnedict\":false,\"dbpedia\":false}},{\"slug\":\"譬え話\",\"is_common\":false,\"tags\":[],\"jlpt\":[],\"japanese\":[{\"word\":\"たとえ話\",\"reading\":\"たとえばなし\"},{\"word\":\"例え話\",\"reading\":\"たとえばなし\"},{\"word\":\"譬え話\",\"reading\":\"たとえばなし\"},{\"word\":\"譬話\",\"reading\":\"たとえばなし\"}],\"senses\":[{\"english_definitions\":[\"allegory\",\"fable\",\"parable\"],\"parts_of_speech\":[\"Noun\"],\"links\":[],\"tags\":[],\"restrictions\":[],\"see_also\":[],\"antonyms\":[],\"source\":[],\"info\":[]},{\"english_definitions\":[\"Parable\"],\"parts_of_speech\":[\"Wikipedia definition\"],\"links\":[{\"text\":\"Read “Parable” on English Wikipedia\",\"url\":\"http://en.wikipedia.org/wiki/Parable?oldid=495216000\"},{\"text\":\"Read “たとえ話” on Japanese Wikipedia\",\"url\":\"http://ja.wikipedia.org/wiki/たとえ話?oldid=41347360\"}],\"tags\":[],\"restrictions\":[],\"see_also\":[],\"antonyms\":[],\"source\":[],\"info\":[],\"sentences\":[]}],\"attribution\":{\"jmdict\":true,\"jmnedict\":false,\"dbpedia\":\"http://dbpedia.org/resource/Parable\"}},{\"slug\":\"518697a3d5dda7b2c604bfe0\",\"tags\":[],\"jlpt\":[],\"japanese\":[{\"reading\":\"たとえばこんなラヴ・ソング\"}],\"senses\":[{\"english_definitions\":[\"Tatoeba Konna Love Song\"],\"parts_of_speech\":[\"Wikipedia definition\"],\"links\":[{\"text\":\"Read “Tatoeba Konna Love Song” on English Wikipedia\",\"url\":\"http://en.wikipedia.org/wiki/Tatoeba_Konna_Love_Song?oldid=484839140\"},{\"text\":\"Read “たとえばこんなラヴ・ソング” on Japanese Wikipedia\",\"url\":\"http://ja.wikipedia.org/wiki/たとえばこんなラヴ・ソング?oldid=42405209\"}],\"tags\":[],\"restrictions\":[],\"see_also\":[],\"antonyms\":[],\"source\":[],\"info\":[],\"sentences\":[]}],\"attribution\":{\"jmdict\":false,\"jmnedict\":false,\"dbpedia\":\"http://dbpedia.org/resource/Tatoeba_Konna_Love_Song\"}}]}"},
	}

	go func() {
		s := NewServer()
		s.Start()
	}()

	time.Sleep(1 * time.Second)

	var wg sync.WaitGroup
	wg.Add(len(tests))

	for _, test := range tests {
		go func(word, expected string) {
			defer wg.Done()

			result, err := queryJisho(word)
			if err != nil {
				t.Errorf("failed to query jisho: %v", err)
			}

			if result != expected {
				t.Errorf("expected %s, got %s", expected, result)
			}
		}(test.Word, test.Expected)
	}

	wg.Wait()
}

// compareMaps compares two maps for equality
func compareMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, valueA := range a {
		if valueB, ok := b[key]; !ok || valueA != valueB {
			return false
		}
	}
	return true
}
