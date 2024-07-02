package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kingster/go-urlshortner/uss"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	appConfig, err := godotenv.Read()
	if err != nil {
		log.Fatalf("Error reading .env file %s", err.Error())
	}

	uss.MetaStore = &uss.USSStore{}
	err = uss.MetaStore.Connect(appConfig)
	if err != nil {
		log.Fatalf("Failed to connect to db %s", err.Error())
	}

	StartHTTPServer()
}

func StartHTTPServer() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${remote_ip} [${time_rfc3339}] "${method} ${uri} HTTP/1.0" ${status} ${latency_human} ${bytes_out} ${error} "${user_agent}"` + "\n",
		Skipper: func(c echo.Context) bool {
			if c.Request().RequestURI == "/healthcheck" {
				return true
			}
			return false
		},
	}))
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "good!")
	})

	e.GET("/resolve/:code", func(c echo.Context) error {
		code := c.Param("code")
		info := uss.MetaStore.FindByShortCode(code)
		if info != nil {
			return c.JSON(http.StatusOK, info)
		} else {
			return c.String(http.StatusNotFound, "Not Found.")
		}
	})

	e.POST("/shorten", func(c echo.Context) error {
		req := new(uss.ShortCodeInfo)
		if err := c.Bind(req); err != nil {
			return err
		}

		req.ShortCode = uss.GenerateShortLink(req.URL)
		err := uss.MetaStore.Persist(req)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed Persisiting Entity with Error %s", err.Error()))
		} else {
			return c.JSON(http.StatusOK, req)
		}
	})

	// automatically add routers for net/http/pprof e.g. /debug/pprof, /debug/pprof/heap, etc.
	// go get github.com/hiko1129/echo-pprof
	//echopprof.Wrap(e)
	e.Logger.Fatal(e.Start(":9090"))
}
