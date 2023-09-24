package main

import (
	"log"
	"net/http"
	"time"

	"phenikaa/infrastructure"
	"phenikaa/router"
)

func main() {

	// go run main.go
	infrastructure.InfoLog.Println("Database name: ", infrastructure.GetDBName())
	log.Printf("Server running at port: %+v\n", infrastructure.GetAppPort())
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router.Router(),
		ReadTimeout:    6000 * time.Second,
		WriteTimeout:   6000 * time.Second,
		MaxHeaderBytes: 1 << 30,
	}
	log.Fatal(s.ListenAndServe())
}
