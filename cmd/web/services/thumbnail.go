package services

import (
	"log"
	"net/http"

	"github.com/linhtutkyawdev/netflixify/cmd/web/components"
)

func ThumbnailWebHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	component := components.Hello(r.Form.Get("name"))
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in CaptureWebHandler: %e", err)
	}
}
