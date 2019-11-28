package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"nquicenob.com/golang-api-example/internal/config"
)

var loggerlvl = map[string]log.Lvl{
	"debug": log.DEBUG,
	"info":  log.INFO,
	"warn":  log.WARN,
	"error": log.ERROR,
	"off":   log.OFF,
}

func New(config *config.Specification) *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(loggerlvl[config.LogLevel])

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	e.GET("/_health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.Validator = NewValidator()
	e.Logger.Debug("Server is ready")
	return e
}
