package router

import (
	"net/http"

	echoSwagger "github.com/drewsilcock/echo-swagger"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/drewsilcock/hbaas-server/docs"
	"github.com/drewsilcock/hbaas-server/handlers"
	"github.com/drewsilcock/hbaas-server/version"
)

func New(db *gorm.DB) *echo.Echo {
	e := echo.New()

	e.Pre(echoMiddleware.RemoveTrailingSlash())

	swaggerHandler := echoSwagger.WrapHandler

	docs.SwaggerInfo.Version = version.Version
	swaggerGroup := e.Group("/")
	swaggerGroup.GET("swagger/*", swaggerHandler)

	apiGroup := e.Group("/")

	handlers.DB = db
	handlers.ConfigureHandlers(apiGroup)

	// If you hit the base URL of the API, you should be redirected to the Swagger UI. Must go after API group so that
	// it takes precedence over them.
	swaggerGroup.GET("", swaggerBaseRedirect)

	return e
}

func swaggerBaseRedirect(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "swagger/index.html")
}
