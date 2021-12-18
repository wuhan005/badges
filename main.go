package main

import (
	"os"
	"strconv"

	"github.com/flamego/flamego"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/badges/internal/context"
	"github.com/wuhan005/badges/internal/route"
)

func main() {
	defer log.Stop()
	if err := log.NewConsole(); err != nil {
		panic(err)
	}

	f := flamego.Classic()
	f.Use(context.Contexter())

	f.Get("/", route.ContributorsBadgeHandler)

	port := 2830
	var err error
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port, err = strconv.Atoi(portEnv)
		if err != nil {
			panic("port must be an integer")
		}
	}

	f.Run(port)
}
