package route

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/disintegration/imaging"

	"github.com/wuhan005/badges/github"
	"github.com/wuhan005/badges/internal/context"
	"github.com/wuhan005/badges/internal/fileutil"
)

func ContributorsBadgeHandler(ctx context.Context) {
	owner := ctx.Param("owner")
	repo := ctx.Param("repo")

	contributors, err := github.GetContributors(owner, repo)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to get contributors: %v", err)
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
		// Download image to temporary path.
		imagePath := path.Join(os.TempDir(), strconv.Itoa(int(contributor.GetID())))
		err := fileutil.DownloadFile(contributor.GetAvatarURL(), imagePath)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to download contributor's avatar: %v", err)
			return
		}

		// Load image file
		img, err := imaging.Open(imagePath)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to load image file: %v", err)
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
	err = encoder.Encode(ctx.ResponseWriter(), background)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to encode image: %v", err)
		return
	}
}
