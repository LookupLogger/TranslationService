package server

import "strings"

// extractSearchTerm extracts the search term from the URL path. eg. /search/hello -> hello
func extractSearchTerm(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return ""
	}

	return strings.TrimSpace(parts[2])
}
