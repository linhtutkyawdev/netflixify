package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/linhtutkyawdev/netflixify/cmd/web"
	"github.com/linhtutkyawdev/netflixify/cmd/web/components"
	"github.com/linhtutkyawdev/netflixify/cmd/web/services"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/", echo.WrapHandler(templ.Handler(components.Index())))
	e.GET("/thumbnail", services.ThumbnailShowHandler)
	// templ.Handler(components.Thumbnail("assets/img/bg.jpeg", "", "", []string{}))
	e.GET("/api", s.ApiHandler)

	e.POST("/thumbnail", echo.WrapHandler(http.HandlerFunc(services.ThumbnailDefaultHandler)))

	return e
}

func (s *Server) ApiHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Api is Working",
	}

	return c.JSON(http.StatusOK, resp)
}
