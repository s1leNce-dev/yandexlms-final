package main

import (
	"final/app"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("[FATAL] %s", err.Error())
	}
}

func main() {
	app.Start()
}
