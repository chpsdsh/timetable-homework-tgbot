package main

import (
	"log"
	"timetable-homework-tgbot/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("key.env not found")
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
