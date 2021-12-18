package main

import (
	"github.com/flamego/flamego"

	"github.com/wuhan005/badges/internal/context"
	"github.com/wuhan005/badges/internal/route"
)

func main() {
	f := flamego.Classic()
	f.Use(context.Contexter())

	f.Get("/repos/{owner}/{repo}/contributors", route.ContributorsBadgeHandler)

	f.NotFound(func(ctx flamego.Context) {
		ctx.Redirect("https://api.github.com/")
	})

	f.Run()
}
