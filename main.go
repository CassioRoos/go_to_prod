package main

import (
	"github.com/CassioRoos/go_to_prod/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/nicholasjackson/env"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var port = env.String("PORT", false, ":8888", "Bind the application to a specific port")

func main() {
	env.Parse()
	echoInstance := echo.New()

	newRelicApp := middleware.ConfigureNewRelic(echoInstance)

	echoInstance.Use(echomiddleware.Logger())
	echoInstance.Use(echomiddleware.Recover())
	echoInstance.Use(middleware.MetricsMiddleware())
	echoInstance.Use(middleware.NewRelicMiddleware(newRelicApp))

	echoInstance.GET("/healthcheck", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "WORKING")
	})
	echoInstance.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	echoInstance.Logger.Fatal(echoInstance.Start(*port))
}
