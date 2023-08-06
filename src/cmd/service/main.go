package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/nats-io/stan.go"
	"gitlab.com/Hofsiedge/l0/internal/config"
)

func main() {
	cfg, err := config.Read[config.ServiceConfig]()
	if err != nil {
		log.Fatal(err)
	}

	opts := []stan.Option{
		stan.NatsURL(cfg.StanURL),
	}
	log.Printf("connecting, %v", cfg)
	sc, err := stan.Connect(cfg.StanCluster, cfg.StanClient, opts...)
	if err != nil {
		err = fmt.Errorf("could not connect to a NATS Streaming server: %w", err)
		log.Fatal(err)
	}
	log.Println("connected to NATS Streaming server")
	defer func() {
		sc.Close()
		log.Println("closed NATS-Streaming connection")
	}()

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		err = fmt.Errorf("could not connect to postgres: %w", err)
		log.Fatal(err)
	}
	log.Println("connected to postgres")
	defer conn.Close(context.Background())

	// TODO: receive all new messages or from last received?
	sub, err := sc.Subscribe(cfg.StanSubject, func(msg *stan.Msg) {
		data := msg.Data
		log.Printf("received a message: %v\n", data)
	}, stan.DeliverAllAvailable())
	if err != nil {
		err = fmt.Errorf("could not subscribe to subject %v: %w", cfg.StanSubject, err)
		log.Fatal(err)
	}
	defer func() {
		if err = sub.Close(); err != nil {
			err = fmt.Errorf("failed subscription closing: %w", err)
			log.Fatal(err)
		} else {
			log.Println("closed subscription")
		}
	}()
	log.Println("subscribed, ready to read")

	// TODO: http handler
	// TODO: http server

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// TODO: uncomment go when http server is implemented
	// go
	func() {
		<-c
		log.Println("shutting down")
		os.Exit(0)
	}()
}
