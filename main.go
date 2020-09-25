package main

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/nicholasjackson/env"
)

var port = env.String("PORT", false, ":8888", "Bind the application to a specific port")

func main() {
	env.Parse()
	echoInstance := echo.New()

	echoInstance.Use(echomiddleware.Logger())
	echoInstance.Use(echomiddleware.Recover())
	echoInstance.Logger.Fatal(echoInstance.Start(*port))
}

