package route

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"path"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"

	"github.com/wuhan005/badges/github"
	"github.com/wuhan005/badges/util"
)

func ContributorsBadgeHandler(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	contributors, err := github.GetContributors(owner, repo)
	if err != nil {
		c.String(500, "Failed to get contributors: %v", err)
		return
	}

	count := len(contributors)
	width := 1280
	imgSize := 96
	padding := 10

	perLine := width / (imgSize + padding)
	line := int(math.Ceil(float64(count) / float64(perLine)))

	background := imaging.New(width, line*(imgSize+padding), color.NRGBA{})

	for index, contributor := range contributors {
		// Download image
		imagePath := path.Join("/tmp", strconv.Itoa(int(contributor.GetID())))
		err := util.DownloadFile(contributor.GetAvatarURL(), imagePath)
		if err != nil {
			c.String(500, "Failed to download contributor's avatar: %v", err)
			return
		}

		// Load image file
		img, err := imaging.Open(imagePath)
		if err != nil {
			c.String(500, "Failed to load image file: %v", err)
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
	err = encoder.Encode(c.Writer, background)
	if err != nil {
		c.String(500, "Failed to encode image: %v", err)
		return
	}
}
