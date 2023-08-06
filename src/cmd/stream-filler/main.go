package main

import (
	"fmt"
	"log"

	"github.com/nats-io/stan.go"
	"gitlab.com/Hofsiedge/l0/internal/config"
)

func main() {
	cfg, err := config.Read[config.FillerConfig]()
	if err != nil {
		log.Fatal(err)
	}

	opts := []stan.Option{
		stan.NatsURL(cfg.StanURL),
	}
	sc, err := stan.Connect(cfg.StanCluster, cfg.StanClient, opts...)
	if err != nil {
		err = fmt.Errorf("could not connect to a NATS Streaming server: %w", err)
		log.Fatal(err)
	}
	log.Println("connected to NATS Streaming server")
	defer func() {
		sc.Close()
		log.Println("closed connection")
	}()
	messages := []string{
		"Message 1",
		"This one is a bit longer",
	}
	log.Println("publishing messages...")
	for _, msg := range messages {
		err = sc.Publish(cfg.StanSubject, []byte(msg))
	}
	log.Println("finished publishing")
}
