package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/repos/:owner/:repo/contributors", contributorsBadgeHandler)

	panic(r.Run())
}
