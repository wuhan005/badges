package main

import (
	"github.com/gin-gonic/gin"

	"github.com/wuhan005/badges/route"
)

func main() {
	r := gin.Default()

	r.GET("/repos/:owner/:repo/contributors", route.ContributorsBadgeHandler)

	r.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "https://api.github.com/")
	})

	panic(r.Run())
}
