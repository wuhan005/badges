package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/repos/:owner/:repo/contributors", contributorsBadgeHandler)
	r.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "https://api.github.com/")
	})

	panic(r.Run())
}
