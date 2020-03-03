package main

import (
	"context"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v29/github"
	"image"
	"image/color"
	"image/png"
	"math"
)

func contributorsBadgeHandler(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	client := github.NewClient(nil)
	ctx := context.Background()

	contributors := make([]*github.Contributor, 0)
	page := 1
	for {
		var tmpContributors []*github.Contributor
		tmpContributors, _, err := client.Repositories.ListContributors(ctx, owner, repo, &github.ListContributorsOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 50,
			},
		})
		if err != nil {
			c.JSON(makeErrJSON(500, 50000, err))
			return
		}
		page++
		if len(tmpContributors) == 0 {
			break
		}
		contributors = append(contributors, tmpContributors...)
	}

	count := len(contributors)
	width := 1280
	imgSize := 96
	padding := 10

	perLine := width / (imgSize + padding)
	line := int(math.Ceil(float64(count) / float64(perLine)))

	background := imaging.New(width, line*(imgSize+padding), color.NRGBA{0, 0, 0, 0})

	for index, contributor := range contributors {
		// Download image
		imgPath, err := downloadFile(contributor.GetAvatarURL(), "/tmp")
		if err != nil {
			c.JSON(makeErrJSON(500, 50000, err))
			return
		}
		// Load image file
		img, err := imaging.Open(imgPath)
		if err != nil {
			c.JSON(makeErrJSON(500, 50000, err))
			return
		}
		// Resize image
		img = imaging.Resize(img, imgSize, imgSize, imaging.Lanczos)
		// Paste image into background
		y := int(math.Floor(float64((index)/perLine))) * (imgSize + padding)
		x := (index - int(math.Floor(float64((index)/perLine)))*perLine) * (imgSize + padding)
		background = imaging.Paste(background, img, image.Pt(x, y))
	}

	encoder := png.Encoder{}
	_ = encoder.Encode(c.Writer, background)
}
