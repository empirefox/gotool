package web

import (
	"mime"
	"net/http"
	"strings"
)

// We decide by the first match
func AcceptJson(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if accept == "" {
		return true
	}
	for _, a := range strings.Split(accept, ",") {
		mediaType, _, _ := mime.ParseMediaType(a)
		switch mediaType {
		case "*/*", "application/*", "application/json":
			return true
		case "text/html", "text/plain":
			return false
		}
	}
	return false
}

func RequestJson(r *http.Request) bool {
	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	return content == "application/json"
}
