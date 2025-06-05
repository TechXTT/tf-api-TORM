package main

import (
	route "github.com/TechXTT/tf-api-TORM/cmd/router"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	r := route.NewRouter()
	r.Run()
}
