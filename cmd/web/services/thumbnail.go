package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/linhtutkyawdev/netflixify/cmd/web/components"

	imgBB "github.com/JohnNON/ImgBB"
)

const (
	key = "3aa950d66034374fe3e87df0f6a1cbc5"
)

func ThumbnailWebHandler(w http.ResponseWriter, r *http.Request) {
	// err := r.ParseForm()

	// if err != nil {
	// 	http.Error(w, "Bad Request", http.StatusBadRequest)
	// }

	file, _, _ := r.FormFile("file")
	defer file.Close()

	// create a destination file
	// dst, _ := os.Create(filepath.Join("./cmd/web/assets/img/", header.Filename))
	// defer dst.Close()

	// upload the file to destination path
	// nb_bytes, _ := io.Copy(dst, file)

	// fmt.Println("File uploaded successfully. ", nb_bytes)

	res := wev(file)
	component := components.Thumbnail(res.Data.Image.URL)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in CaptureWebHandler: %e", err)
	}
}

// const domain = "http://localhost:8080"
// const defaultUsername = "doctron"
// const defaultPassword = "lampnick"

// func createThumbnail() {
// 	client := doctron.NewClient(context.Background(), domain, defaultUsername, defaultPassword)
// 	req := doctron.NewDefaultHTML2ImageRequestDTO()
// 	req.ConvertURL = "http://localhost:3000/thumbnail"
// 	response, err := client.HTML2Image(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println(len(response.Data))

// 	f, err := os.Create("example.png")
// 	if err != nil {
// 		log.Fatal("Cannot create file", err)
// 	}
// 	defer f.Close()

// 	n2, err := f.Write(response.Data)
// 	if err != nil {
// 		log.Fatal("Cannot write file", err)
// 	}
// 	fmt.Printf("wrote %d bytes\n", n2)

// 	f2, err := os.Open("example.png")
// 	if err != nil {
// 		log.Fatal("Cannot open file", err)
// 	}

// 	img, err := png.Decode(f2)
// 	if err != nil {
// 		log.Fatal("Cannot decode image:", err)
// 	}

// 	cImg, err := cutter.Crop(img, cutter.Config{
// 		Height: 405, // height in pixel or Y ratio(see Ratio Option below)
// 		Width:  720, // width in pixel or X ratio
// 		// Mode:    cutter.TopLeft,      // Accepted Mode: TopLeft, Centered
// 		// Anchor:  image.Point{10, 10}, // Position of the top left point
// 		// Options: 0,                   // Accepted Option: Ratio
// 	})

// 	if err != nil {
// 		log.Fatal("Cannot crop image:", err)
// 	}

// 	f3, err := os.Create("example.png")
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer f3.Close()
// 	if err := png.Encode(f3, cImg); err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("cImg dimension:", cImg.Bounds())

// }

func wev(f multipart.File) imgBB.Response {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		log.Fatal(err)
	}
	img, err := imgBB.NewImageFromFile(hashSum(buf.Bytes()), 7*24*60*60, buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	imgBBClient := imgBB.NewClient(httpClient, key)

	resp, err := imgBBClient.Upload(context.Background(), img)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", resp)
	return resp
}

func hashSum(b []byte) string {
	sum := md5.Sum(b)

	return hex.EncodeToString(sum[:])
}
