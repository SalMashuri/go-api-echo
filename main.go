package main

import (
	"fmt"
	"os"
	"sync"

	localConfig "github.com/Rifannurmuhammad/go-api-echo/config"
	"github.com/Rifannurmuhammad/go-api-echo/config/rsa"
	config "github.com/joho/godotenv"
)

func main() {
	err := config.Load(".env")
	if err != nil {
		fmt.Println(".env is not loaded properly")
		os.Exit(2)
	}

	// initiate database and other connections
	readDB := localConfig.ReadPostgresDB()
	writeDB := localConfig.WritePostgresDB()

	service := MakeHandler(readDB, writeDB)

	wg := sync.WaitGroup{}

	wg.Add(1)
	publicKey, err := rsa.InitPublicKey()
	go func() {
		defer wg.Done()
		service.HTTPServerMain(publicKey)
	}()
	// Wait All services to end
	wg.Wait()
}
