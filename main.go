package main

import (
	"os"
	"strconv"

	"github.com/flamego/flamego"

	"github.com/wuhan005/badges/internal/context"
	"github.com/wuhan005/badges/internal/route"
)

func main() {
	f := flamego.Classic()
	f.Use(context.Contexter())

	f.Get("/", route.ContributorsBadgeHandler)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("port must be an integer")
	}
	f.Run(port)
}
