package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/stan.go"
	"gitlab.com/Hofsiedge/l0/internal/config"
	"gitlab.com/Hofsiedge/l0/internal/repo/postgres"
	"gitlab.com/Hofsiedge/l0/internal/server"
)

func main() {
	cfg, err := config.Read[config.ServiceConfig]()
	if err != nil {
		log.Fatal(err)
	}

	// stan connection
	sc, err := stan.Connect(cfg.StanCluster, cfg.StanClient, stan.NatsURL(cfg.StanURL))
	if err != nil {
		err = fmt.Errorf("could not connect to a NATS Streaming server: %w", err)
		log.Fatal(err)
	}
	log.Println("connected to NATS-Streaming server")
	defer func() {
		sc.Close()
		log.Println("closed NATS-Streaming connection")
	}()

	// postgresql connection
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		err = fmt.Errorf("could not connect to postgres: %w", err)
		log.Fatal(err)
	}
	log.Println("connected to postgres")
	defer conn.Close(context.Background())

	srv, err := server.NewServer(new(postgres.MockOrders))
	if err != nil {
		err = fmt.Errorf("failed to create a server: %w", err)
		log.Fatal(err)
	}

	// subscription
	srv.Stan, err = sc.Subscribe(cfg.StanSubject, srv.HandleMessage, stan.DurableName(cfg.StanDurableName))
	if err != nil {
		err = fmt.Errorf("could not subscribe to subject %v: %w", cfg.StanSubject, err)
		log.Fatal(err)
	}
	defer func() {
		if err = srv.Stan.Close(); err != nil {
			err = fmt.Errorf("failed subscription closing: %w", err)
			log.Fatal(err)
		} else {
			log.Println("closed subscription")
		}
	}()
	log.Println("subscribed, ready to read")

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1/order").Subrouter()
	s.HandleFunc("/", srv.ListEndpoint).Methods("GET")
	s.HandleFunc("/{id}", srv.GetByIdEndpoint).Methods("GET")

	log.Fatal(http.ListenAndServe("0.0.0.0:80", r))
}
