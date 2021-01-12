package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/awoodbeck/gnp/ch14/gcp"
)

func main() {
	port, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		log.Fatal("FUNCTIONS_CUSTOMHANDLER_PORT environment variable not set")
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           http.HandlerFunc(gcp.LatestXKCD),
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
	}

	log.Printf("Listening on %q ...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
