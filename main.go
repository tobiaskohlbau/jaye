// Copyright (c) 2017 Tobias Kohlbau
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"kohlbau.de/x/jaye/config"
	"kohlbau.de/x/jaye/handler"
	"kohlbau.de/x/jaye/services/youtube"
)

func main() {
	configPath := flag.String("config", "./config/config.json", "path to config file")

	flag.Parse()

	// Read config
	config, err := config.FromFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// YouTube Service
	ytService := youtube.New(config.Youtube.URL, config.Youtube.Token, config.Youtube.VideoPath)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: handler.New(ytService),
	}

	go func() {
		// Graceful shutdown
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, os.Kill)

		<-sigquit

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Unable to shut down server: %v", err)
			return
		}
	}()

	// Start server
	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
		return
	}
}
