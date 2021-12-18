package main

import (
	"os"

	"github.com/flamego/flamego"

	"github.com/wuhan005/badges/internal/context"
	"github.com/wuhan005/badges/internal/route"
)

func main() {
	f := flamego.Classic()
	f.Use(context.Contexter())

	f.Get("/", route.ContributorsBadgeHandler)

	port := os.Getenv("PORT")
	f.Run(port)
}
