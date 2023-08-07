package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/nats-io/stan.go"
	"gitlab.com/Hofsiedge/l0/internal/config"
	"gitlab.com/Hofsiedge/l0/internal/domain"
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
	log.Println("connected to NATS-Streaming server")
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

	sub, err := sc.Subscribe(cfg.StanSubject, func(msg *stan.Msg) {
		data := msg.Data
		var order domain.Order
		if err := json.Unmarshal(data, &order); err != nil {
			log.Printf("could not unmarshal a message: %v", data)
			return
		}
		log.Printf("received a message with an Order: %v\n", order)
	}, stan.DurableName(cfg.StanDurableName))
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
