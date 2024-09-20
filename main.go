package main

import (
	"bufio"
	"context"
	"fmt"
	"nats_bucket_walker/cli"
	"os"
	"time"

	"github.com/buger/goterm"
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

	// DEPRECATED
	keys, err := kv.Keys(ctx)
	if err != nil {
		return []string{}, err
	}

	return keys, nil
}
func GetAllBuckets() ([]string, error) {
	// connect to nats server
	url := nats.DefaultURL
	if val, ok := os.LookupEnv("NATS_URL"); ok {
		url = val
	} else {
		panic("Could not connect to NATS")
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
		println(n)
		names = append(names, n)
	}

	return names, nil
}

func PressToContinue() {
	fmt.Print("press to continue...")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())
}

func main() {
	buckets, err := GetAllBuckets()
	if err != nil {
		panic(err)
	}

	for _, k := range buckets {
		fmt.Println(k)
	}

	server := "#SERVER"
	for {
		goterm.Clear()
		menu := cli.NewMenu(fmt.Sprintf("%v Available buckets", server))

		for _, b := range buckets {
			menu.AddItem(b, b)

		}
		menu.AddItem("🚪 Quit", "quit")

		choice := menu.Display()

		if choice == "quit" {
			break
		} else {
			// list all keys in a bucket
			keys, err := GetAllKeys(choice)
			if err != nil {
				panic(err)
			}
			fmt.Println("==================")
			fmt.Printf(" Content of %v\n", choice)
			fmt.Println("==================")
			for _, k := range keys {
				println(k)
			}
			println("==================")
			PressToContinue()

		}

	}

}
