package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/linhtutkyawdev/netflixify/cmd/web/components"
)

func ThumbnailWebHandler(w http.ResponseWriter, r *http.Request) {
	// err := r.ParseForm()

	// if err != nil {
	// 	http.Error(w, "Bad Request", http.StatusBadRequest)
	// }

	file, header, _ := r.FormFile("file")
	defer file.Close()

	// create a destination file
	dst, _ := os.Create(filepath.Join("./", header.Filename))
	defer dst.Close()

	// upload the file to destination path
	nb_bytes, _ := io.Copy(dst, file)

	fmt.Println("File uploaded successfully. ", nb_bytes)

	component := components.Hello("")
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in CaptureWebHandler: %e", err)
	}
}
