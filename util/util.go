package util

import (
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func DownloadFile(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "http GET")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "open file")
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Wrap(err, "copy")
	}

	return nil
}
