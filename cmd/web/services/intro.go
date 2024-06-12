package services

import (
	"github.com/labstack/echo/v4"
	"github.com/linhtutkyawdev/netflixify/cmd/web/components"
)

func Intro(c echo.Context) error {
	href := c.QueryParam("href")
	component := components.LetterAnimation(href)
	return component.Render(c.Request().Context(), c.Response())
}
