package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gre-ory/games-go/internal/game/tictactoe"
)

// //////////////////////////////////////////////////
// main

func main() {

	// exit process immediately upon sigterm
	handleSigTerms()

	// logger
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	// router
	router := http.NewServeMux()

	// serve static files
	registerStaticFiles(router)

	// register games
	tictactoe.Register(router)

	// listen
	port := "9090"
	logger.Println("listening on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Println("http.ListenAndServe():", err)
		os.Exit(1)
	}
}

// //////////////////////////////////////////////////
// static

var (
	//go:embed static/*
	staticFS embed.FS
)

func registerStaticFiles(router *http.ServeMux) {
	router.Handle("/static/", http.FileServer(http.FS(staticFS)))
}

// //////////////////////////////////////////////////
// cors

// cors := func(h http.Handler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// in development, the Origin is the the Hugo server, i.e. http://localhost:1313
// 		// but in production, it is the domain name where one's site is deployed
// 		//
// 		// CHANGE THIS: You likely do not want to allow any origin (*) in production. The value should be the base URL of
// 		// where your static content is served
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, hx-target, hx-current-url, hx-request")
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
// 		h.ServeHTTP(w, r)
// 	}
// }

// //////////////////////////////////////////////////
// sigterms

func handleSigTerms() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")
		os.Exit(1)
	}()
}
