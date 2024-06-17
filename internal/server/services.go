package server

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	imgBB "github.com/JohnNON/ImgBB"
	"github.com/chromedp/chromedp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"github.com/linhtutkyawdev/netflixify/cmd/web/components"
)

const (
	imgBB_key     = "3aa950d66034374fe3e87df0f6a1cbc5"
	defaultImgSrc = "/assets/images/thumbnail-bg.jpeg"
)

func ApiHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Api is Working",
	}

	return c.JSON(http.StatusOK, resp)
}

func IntroHandler(c echo.Context) error {
	href := c.QueryParam("href")
	component := components.LetterAnimation(href)
	return component.Render(c.Request().Context(), c.Response())
}

func (s *Server) ThumbnailPostHandler(w http.ResponseWriter, r *http.Request) {
	imgSrc := ""
	file, _, _ := r.FormFile("imgFile")
	if file != nil {
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			log.Fatal(err)
		}
		imgSrc = uploadToImgbb(buf.Bytes(), 5*60).Data.Image.URL
		defer file.Close()
	} else {
		imgSrc = r.FormValue("imgSrc")
	}
	originalUrl := "/thumbnail?imgSrc=" + imgSrc + "&title=" + r.FormValue("title") + "&subtitle=" + r.FormValue("subtitle") + "&categories=" + r.FormValue("categories")

	thumbnailUrl := screenshotAndUplload(os.Getenv("URL")+originalUrl, `#thumbnail-container`).Data.Image.URL

	if r.FormValue("id") != "" {
		id, _ := strconv.Atoi(r.FormValue("id"))

		s.bot.Send(tgbotapi.NewPhoto(int64(id), tgbotapi.FileURL(thumbnailUrl)))
		component := components.CloseWebApp()

		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Fatalf("Error rendering in CaptureWebHandler: %e", err)
		}
	} else {
		component := components.DownloadThumbnail(thumbnailUrl, strings.ToLower(strings.ReplaceAll(r.FormValue("title"), " ", "")), originalUrl)

		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Fatalf("Error rendering in CaptureWebHandler: %e", err)
		}
	}
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
	id := c.QueryParam("tgId")

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
	component := components.Thumbnail(id, imgSrc, title, subtitle, categories)
	return component.Render(c.Request().Context(), c.Response())
}

func screenshotAndUplload(url string, sel string) imgBB.Response {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(url, sel, &buf)); err != nil {
		log.Fatal(err)
	}

	return uploadToImgbb(buf, 60*60)
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
