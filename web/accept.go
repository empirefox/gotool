package web

import (
	"mime"
	"net/http"
	"strings"
)

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
		}
	}
	return false
}
