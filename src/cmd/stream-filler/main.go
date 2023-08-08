package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	_ "embed"

	"github.com/nats-io/stan.go"
	"gitlab.com/Hofsiedge/l0/internal/config"
)

// read messages from cfg.DataDir. skips files on errors.
func readMessages(cfg config.FillerConfig) ([][]byte, error) {
	messages := make([][]byte, 0)
	entries, err := os.ReadDir(cfg.DataDir)
	if err != nil {
		return nil, fmt.Errorf("could not read the messages: %w", err)
	}
	for _, entry := range entries {
		filePath := path.Join(cfg.DataDir, entry.Name())
		if entry.IsDir() || !strings.HasSuffix(filePath, ".json") {
			log.Printf("skipping %v", filePath)
			continue
		}
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("could not read %s, skipping", filePath)
			continue
		}
		message, err := io.ReadAll(file)
		if err != nil {
			log.Printf("IO error reading %v, skipping", err)
			continue
		}
		messages = append(messages, message)
	}
	return messages, nil
}

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

	messages, err := readMessages(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("waiting before publishing")

	time.Sleep(20 * time.Second)

	log.Println("publishing messages...")
	for _, msg := range messages {
		err = sc.Publish(cfg.StanSubject, []byte(msg))
	}
	log.Println("finished publishing")
}
