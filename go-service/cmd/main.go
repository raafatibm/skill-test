package main

import (
	"log"
	"os"
	"os/signal"
	"reporting/internal/models"
	"reporting/internal/server"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("failed to loading .env file with error: ", err)
		os.Exit(1)
	}

	cfg, err := models.GetConfig()
	if err != nil {
		log.Fatal("failed to retrieve configuration with error: ", err)
		os.Exit(2)
	}

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal("failed to start the service server with error: ", err)
		os.Exit(3)
	}

	go server.Start()

	sigStop := make(chan os.Signal, 1)
	signal.Notify(sigStop, os.Interrupt)
	signal.Notify(sigStop, syscall.SIGTERM)

	stopped := <-sigStop
	log.Println(stopped.String() + " signal received")

	err = server.Shutdown()
	if err != nil {
		log.Println("failed to shutdown the http server with error: ", err)
		os.Exit(4)
	}

	log.Println("exiting reporting service")
}
