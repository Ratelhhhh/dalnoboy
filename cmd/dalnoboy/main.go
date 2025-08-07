package main

import (
	"log"
	"os"

	"dalnoboy/internal/app"
)

func main() {
	appName := "dalnoboy"
	if len(os.Args) > 1 {
		appName = os.Args[1]
	}

	application := app.New(appName)
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
