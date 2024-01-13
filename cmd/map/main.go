package main

import (
	"app/internal/application"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// env
	// ...
	if err := godotenv.Load("./../../.env"); err != nil {
		fmt.Println(err)
		return
	}

	// application
	// - config
	cfg := &application.ConfigAppMap{
		Addr:       os.Getenv("SERVER_ADDR"),
		Token:      os.Getenv("API_TOKEN"),
		LayoutDate: os.Getenv("LAYOUT_DATE"),
	}
	app := application.NewApplicationMap(cfg)
	// - set-up
	if err := app.SetUp(); err != nil {
		fmt.Println(err)
		return
	}

	// run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
