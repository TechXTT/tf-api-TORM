package main

import (
	route "github.com/hacktues-9/tf-api/cmd/router"
	models "github.com/hacktues-9/tf-api/models"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// DB := db.Init()
	client := models.NewClient()
	r := route.NewRouter(client)
	r.Init()
	r.Run()
}
