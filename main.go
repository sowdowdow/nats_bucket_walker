package main

import (
	"context"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func GetAllKeys(bucket string) ([]string, error) {
	// connect to nats server
	url := nats.DefaultURL
	if val, ok := os.LookupEnv("NATS_URL"); ok {
		url = val
	} else {
		panic("Could not connect to NATS")
	}
	conn, _ := nats.Connect(url)
	// defer conn.Close()

	// create jetstream context from nats connection
	js, err := jetstream.New(conn)
	if err != nil {
		return []string{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return []string{}, err
	}

	keys, err := kv.Keys(ctx)
	if err != nil {
		return []string{}, err
	}

	return keys, nil
}

func main() {
	GetAllKeys("example")
}
