package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lampnick/doctron-client-go"
	"github.com/linhtutkyawdev/netflixify/cmd/web/components"

	imgBB "github.com/JohnNON/ImgBB"
)

const (
	// imgBB
	imgBB_key = "3aa950d66034374fe3e87df0f6a1cbc5"

	// doctron
	// https://doctron-latest.onrender.com
	host = "http://0.0.0.0"

	// wev
	temp_file     = "tmp.png"
	defaultImgSrc = "assets/img/bg.jpeg"
)

func ThumbnailPostHandler(w http.ResponseWriter, r *http.Request) {
	imgSrc := ""
	file, _, _ := r.FormFile("imgFile")
	if file != nil {
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			log.Fatal(err)
		}
		imgSrc = uploadToImgbb(buf.Bytes(), 60).Data.Image.URL
		defer file.Close()
	} else {
		imgSrc = r.FormValue("imgSrc")
	}
	originalUrl := "/thumbnail?imgSrc=" + imgSrc + "&title=" + r.FormValue("title") + "&subtitle=" + r.FormValue("subtitle") + "&categories=" + r.FormValue("categories")
	thumbnailUrl := createThumbnail(host + ":" + os.Getenv("PORT") + originalUrl).Data.Image.URL

	component := components.DownloadThumbnail(thumbnailUrl, strings.ToLower(strings.ReplaceAll(r.FormValue("title"), " ", "")), originalUrl)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in CaptureWebHandler: %e", err)
	}
}

func createThumbnail(url string) imgBB.Response {
	client := doctron.NewClient(context.Background(), host+":"+os.Getenv("DOCTRON_PORT"), os.Getenv("DOCTRON_USERNAME"), os.Getenv("DOCTRON_PASSWORD"))
	req := doctron.NewDefaultHTML2ImageRequestDTO()
	req.ConvertURL = url
	req.CustomClip = true
	req.ClipX = 0
	req.ClipY = 0
	req.ClipWidth = 720
	req.ClipHeight = 405
	req.ClipScale = 1.0

	// url to img
	response, err := client.HTML2Image(req)
	if err != nil {
		log.Fatal(err)
	}

	res := uploadToImgbb(response.Data, 60*60)

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
