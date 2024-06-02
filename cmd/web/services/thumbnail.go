package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lampnick/doctron-client-go"
	"github.com/linhtutkyawdev/netflixify/cmd/web/components"
	"github.com/oliamb/cutter"

	imgBB "github.com/JohnNON/ImgBB"
)

const (
	key         = "3aa950d66034374fe3e87df0f6a1cbc5"
	tempImgFile = "tmp.png"
	bgImgUrl    = "assets/img/bg.jpeg"
	Url         = "http://localhost:3000"

	//doctron
	domain          = "http://localhost:8080"
	defaultUsername = "doctron"
	defaultPassword = "lampnick"
)

func ThumbnailDefaultHandler(w http.ResponseWriter, r *http.Request) {
	file, _, _ := r.FormFile("imgSrc")
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Fatal(err)
	}

	bg := upload_to_imgbb(buf.Bytes(), 60)

	thumbnail := createThumbnail(Url + "/thumbnail?imgSrc=" + bg.Data.Image.URL + "&title=" + r.FormValue("title") + "&subtitle=" + r.FormValue("subtitle") + "&categories=" + r.FormValue("categories"))

	w.Header().Set("Content-Type", "image/png")
	component := components.Thumbnail(thumbnail.Data.Image.URL, r.FormValue("title"), r.FormValue("subtitle"), strings.Split(r.FormValue("categories"), ","))
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in CaptureWebHandler: %e", err)
	}
}

func createThumbnail(url string) imgBB.Response {
	client := doctron.NewClient(context.Background(), domain, defaultUsername, defaultPassword)
	req := doctron.NewDefaultHTML2ImageRequestDTO()
	req.ConvertURL = url
	response, err := client.HTML2Image(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(response.Data))

	f, err := os.Create(tempImgFile)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer f.Close()

	n2, err := f.Write(response.Data)
	if err != nil {
		log.Fatal("Cannot write file", err)
	}
	fmt.Printf("wrote %d bytes\n", n2)

	f2, err := os.Open(tempImgFile)
	if err != nil {
		log.Fatal("Cannot open file", err)
	}

	img, err := png.Decode(f2)
	if err != nil {
		log.Fatal("Cannot decode image:", err)
	}

	cImg, err := cutter.Crop(img, cutter.Config{
		Height: 405, // height in pixel or Y ratio(see Ratio Option below)
		Width:  720, // width in pixel or X ratio
		// Mode:    cutter.TopLeft,      // Accepted Mode: TopLeft, Centered
		// Anchor:  image.Point{10, 10}, // Position of the top left point
		// Options: 0,                   // Accepted Option: Ratio
	})

	if err != nil {
		log.Fatal("Cannot crop image:", err)
	}

	buf := bytes.NewBuffer(nil)
	err = png.Encode(buf, cImg)

	if err != nil {
		log.Fatal("Cannot encode image:", err)
	}

	res := upload_to_imgbb(buf.Bytes(), 60*60)

	os.Remove(tempImgFile)
	return res
}

func upload_to_imgbb(bytes []byte, exp uint64) imgBB.Response {
	img, err := imgBB.NewImageFromFile(hashSum(bytes), exp, bytes)
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

func ThumbnailShowHandler(c echo.Context) error {
	imgSrc := c.QueryParam("imgSrc")
	if imgSrc == "" {
		imgSrc = bgImgUrl
	}
	title := c.QueryParam("title")
	if title == "" {
		title = "One Piece"
	}
	subtitle := c.QueryParam("subtitle")
	if subtitle == "" {
		subtitle = "Rating : 9.5"
	}
	categoriesStr := c.QueryParam("categories")
	categories := []string{"Anime", "Series", "Shonen", "Action", "Comedy"}
	if categoriesStr != "" {
		categories = strings.Split(categoriesStr, ",")
	}
	component := components.Thumbnail(imgSrc, title, subtitle, categories)
	return component.Render(c.Request().Context(), c.Response())
}

// ?title=One%20Piece&subtitle=Rating+%3A+9.5&categories=Anime,Series,Shounen,Action,Comedy&imgSrc=assets%2Fimg%2Fbg.jpeg
