package nats

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type WatchUpdate struct {
	Operation Operation
	Key       string
	Value     string
}

type Operation struct {
	Name string
}

func GetAllKV(bucket string) ([]string, error) {
	// connect to nats server
	url := nats.DefaultURL
	if val, ok := os.LookupEnv("NATS_URL"); ok {
		url = val
	} else {
		panic("Could not connect to NATS, no NATS_URL specified")
	}
	conn, _ := nats.Connect(url)
	defer conn.Close()

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

	// DEPRECATED
	keys, err := kv.Keys(ctx)
	kv_map := []string{}
	for _, k := range keys {
		kvEntry, err := kv.Get(ctx, k)
		if err != nil {
			return []string{}, err
		}
		val := string(kvEntry.Value())
		kv_map = append(kv_map, k+" = "+val)
	}
	if err != nil {
		return []string{}, err
	}

	return kv_map, nil
}

func WatchBucket(bucket string, c chan WatchUpdate, quit chan int) error {
	// connect to nats server
	url := nats.DefaultURL
	if val, ok := os.LookupEnv("NATS_URL"); ok {
		url = val
	} else {
		panic("Could not connect to NATS, no NATS_URL specified")
	}
	conn, _ := nats.Connect(url)
	defer conn.Close()

	// create jetstream context from nats connection
	js, err := jetstream.New(conn)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return err
	}

	watcher, err := kv.WatchAll(ctx)
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for u := range watcher.Updates() {
		select {
		case <-quit:
			fmt.Println("quitting bucket watch")
			return nil
		default:
			if u != nil {
				k := string(u.Key())
				v := string(u.Value())

				switch u.Operation() {
				case jetstream.KeyValueDelete:
					c <- WatchUpdate{Operation: Operation{"DEL"}, Key: k}
				case jetstream.KeyValuePut:
					c <- WatchUpdate{Operation: Operation{"PUT"}, Key: k, Value: v}
				case jetstream.KeyValuePurge:
					c <- WatchUpdate{Operation: Operation{"PURGE"}, Key: k}
				default:
					c <- WatchUpdate{Operation: Operation{"UNKNOWN"}, Value: u.Operation().String()}
				}
			}
		}
	}

	return nil
}

func GetAllBuckets() ([]string, error) {
	// connect to nats server
	url := nats.DefaultURL
	if val, ok := os.LookupEnv("NATS_URL"); ok {
		url = val
	} else {
		panic("Could not connect to NATS, no NATS_URL specified")
	}
	conn, _ := nats.Connect(url)
	defer conn.Close()

	// create jetstream context from nats connection
	js, err := jetstream.New(conn)
	if err != nil {
		return []string{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// retrieve all names
	names := []string{}
	osnl := js.KeyValueStoreNames(ctx)
	for n := range osnl.Name() {
		names = append(names, n)
	}

	return names, nil
}
