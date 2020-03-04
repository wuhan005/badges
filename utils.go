package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func makeErrJSON(httpStatusCode int, errCode int, msg interface{}) (int, interface{}) {
	return httpStatusCode, gin.H{"error": errCode, "msg": fmt.Sprint(msg)}
}

func makeSuccessJSON(data interface{}) (int, interface{}) {
	return 200, gin.H{"error": 0, "msg": "success", "data": data}
}

func downloadFile(url string, path string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	path = filepath.Join(path, filepath.Base(url))
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, res.Body)
	defer file.Close()
	if err != nil {
		return "", err
	}

	return path, nil
}
