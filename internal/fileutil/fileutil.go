// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package fileutil

import (
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// DownloadFile downloads file from the given url to the given path.
func DownloadFile(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "http GET")
	}
	defer func() { _ = resp.Body.Close() }()

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Wrap(err, "copy")
	}
	return nil
}
