// Package main is the starting point of the app
package main

import (
	"log"
	"net/http"
	"github.com/cbarlas/PeakAPI/router"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// creating the router
	rt := router.CreateRouter()

	// router listens the port "8081" in "localhost"
	log.Fatal(http.ListenAndServe(":8081", rt))
}

