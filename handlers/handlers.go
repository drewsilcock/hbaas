package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func ConfigureHandlers(g *echo.Group) {
	configureBirthdayRouter(g)
	configurePersonRouter(g)
}

func readPathId(c echo.Context) (int, error) {
	var id int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return id, echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("Invalid format for ID: '%s'.", c.Param("id")),
		)
	}

	return id, nil
}

func urlDecodeParam(c echo.Context, param string) (string, error) {
	uri, err := url.Parse(c.Param(param))
	if err != nil {
		return "", err
	}

	return uri.Path, nil
}
