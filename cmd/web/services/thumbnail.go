package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"image/png"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lampnick/doctron-client-go"
	"github.com/linhtutkyawdev/netflixify/cmd/web/components"
	"github.com/oliamb/cutter"

	imgBB "github.com/JohnNON/ImgBB"
)

const (
	// imgBB
	imgBB_key = "3aa950d66034374fe3e87df0f6a1cbc5"

	// doctron
	doctron_host    = "http://localhost:8080"
	doctronUsername = "doctron"
	doctronPassword = "lampnick"

	// wev
	temp_file     = "tmp.png"
	defaultImgSrc = "assets/img/bg.jpeg"
	host          = "http://localhost:3000"
)

func ThumbnailPostHandler(w http.ResponseWriter, r *http.Request) {
	imgSrc := ""
	file, _, _ := r.FormFile("imgSrc")
	if file != nil {
		println("\n\n\n\nNotNil!\n\n\n\n\n\n")
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			log.Fatal(err)
		}
		imgSrc = uploadToImgbb(buf.Bytes(), 60).Data.Image.URL
		defer file.Close()
	}

	thumbnailUrl := createThumbnail(host + "/thumbnail?imgSrc=" + imgSrc + "&title=" + r.FormValue("title") + "&subtitle=" + r.FormValue("subtitle") + "&categories=" + r.FormValue("categories")).Data.Image.URL

	res, err := http.Get(thumbnailUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, res.Body); err != nil {
		log.Fatal(err)
	}

	component := components.DownloadThumbnail("data:image/png;base64, "+base64.StdEncoding.EncodeToString(buf.Bytes()), strings.ToLower(strings.ReplaceAll(r.FormValue("title"), " ", "_")))
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in CaptureWebHandler: %e", err)
	}
}

func createThumbnail(url string) imgBB.Response {
	client := doctron.NewClient(context.Background(), doctron_host, doctronUsername, doctronPassword)
	req := doctron.NewDefaultHTML2ImageRequestDTO()
	req.ConvertURL = url

	// url to img
	response, err := client.HTML2Image(req)
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(bytes.NewReader(response.Data))

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

	res := uploadToImgbb(buf.Bytes(), 60*60)

	return res
}

func uploadToImgbb(bytes []byte, exp uint64) imgBB.Response {
	img, err := imgBB.NewImageFromFile(hashSum(bytes), exp, bytes)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	imgBBClient := imgBB.NewClient(httpClient, imgBB_key)

	res, err := imgBBClient.Upload(context.Background(), img)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func hashSum(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

func ThumbnailHandler(c echo.Context) error {
	imgSrc := c.QueryParam("imgSrc")
	if imgSrc == "" {
		imgSrc = defaultImgSrc
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
