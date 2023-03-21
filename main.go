package main

import (
	route "github.com/hacktues-9/tf-api/cmd/router"
	db "github.com/hacktues-9/tf-api/pkg/database"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	DB := db.Init()
	r := route.NewRouter(DB)
	r.Init()
	r.Run()
}
