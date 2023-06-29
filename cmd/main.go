package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/VicOsewe/go-shell-command-api/configs"
	"github.com/VicOsewe/go-shell-command-api/presentation"
)

const (
	waitSeconds = 30
)

func main() {
	port, err := strconv.Atoi(
		configs.MustGetEnvVar("PORT"),
	)
	if err != nil {
		log.Fatalln("failed to retrieve port value", err)
	}
	srv := presentation.PrepareServer(port)
	log.Printf("server started successfully")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln("There's an error with the server,", err)
		}
	}()

	// Block until we receive a sigint (CTRL+C) signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*waitSeconds,
	)
	defer cancel()

	//Doesn't block if no connections, but will otherwise wait until timeout
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Printf("error during clean shutdown: %s", err)
		os.Exit(-1)
	}

	log.Printf(
		"graceful shutdown started; the timeout is %d secs",
		waitSeconds,
	)
	os.Exit(0)
}
