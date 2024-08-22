package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/linhtutkyawdev/netflixify/cmd/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/", WebAppHandler)
	e.GET("/api", ApiHandler)
	e.GET("/video", s.VideoHandler)

	e.GET("/thumbnail", ThumbnailHandler)
	e.POST("/thumbnail", echo.WrapHandler(http.HandlerFunc(s.ThumbnailPostHandler)))

	e.GET("/animation", AnimationHandler)

	return e
}
