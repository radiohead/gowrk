package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"radiohead/gowrk/cmd/gowrk/config"
	"radiohead/gowrk/pkg/http"
	"radiohead/gowrk/pkg/runner"
)

func main() {
	log.SetFlags(0)
	config, usage, err := config.New()
	if err != nil {
		log.Printf("error: %s", err.Error())
		log.Println(usage)

		os.Exit(1)
	}

	log.Printf("starting with config %+v\n", config)

	client, err := http.NewPooledClient(config.PoolSize, config.MaxConns)
	if err != nil {
		log.Println("failed to create HTTP client")
		log.Println(err.Error())

		os.Exit(1)
	}

	request, err := http.NewGETRequest(config.URL)
	if err != nil {
		log.Println("failed to create HTTP request")
		log.Println(err.Error())

		os.Exit(1)
	}

	closeCh := make(chan interface{}, 1)
	runner, err := runner.New(config.Rate, config.IsVerbose, request, client, closeCh)
	if err != nil {
		log.Println("failed to create runner")
		log.Println(err.Error())

		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		runner.Start()
	}()

	<-stop
	log.Println("stopping...")

	close(closeCh)
	log.Println("exit")
	os.Exit(0)
}
