package main

import (
	"crypto/rsa"
	"fmt"
	"os"
	"strconv"

	"github.com/Rifannurmuhammad/go-api-echo/config"
	"github.com/Rifannurmuhammad/go-api-echo/middleware"
	memberDelivery "github.com/Rifannurmuhammad/go-api-echo/src/member/delivery"
	"github.com/labstack/echo"
	mid "github.com/labstack/echo/middleware"
)

const (
	//DefaultPort default for http port
	DefaultPort = 8080
)

// HTTPServerMain main function for serving services over http
func (s *Service) HTTPServerMain(publicKey *rsa.PublicKey) {
	e := echo.New()

	e.Use(middleware.ServerHeader, middleware.Logger)
	//e.Use(mid.Recover())
	e.Use(mid.CORS())

	if os.Getenv("DEVELOPMENT") == "1" {
		e.Debug = true
	}

	redisConnection, _ := config.GetRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_TLS"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_PORT"))
	cl := redisConnection

	// member endpoints
	memberHandler := memberDelivery.NewHTTPHandler(s.MemberUseCase)
	meGroup := e.Group("/api/me")
	meGroup.Use(middleware.BearerVerify(publicKey, cl, true))
	memberHandler.MountMe(meGroup)

	// set REST port
	var port uint16
	if portEnv, ok := os.LookupEnv("PORT"); ok {
		portInt, err := strconv.Atoi(portEnv)
		if err != nil {
			port = DefaultPort
		} else {
			port = uint16(portInt)
		}
	} else {
		port = DefaultPort
	}

	listenerPort := fmt.Sprintf(":%d", port)
	e.Logger.Fatal(e.Start(listenerPort))
}
