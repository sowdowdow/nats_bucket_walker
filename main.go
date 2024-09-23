package main

import (
	"context"
	"fmt"
	"nats_bucket_walker/cli"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/buger/goterm"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func Clear() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndexAny(s[:max], " .,:;-")] + "..."
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
	kv_map := []string{}
	for _, k := range keys {
		kvEntry, err := kv.Get(ctx, k)
		if err != nil {
			return []string{}, err
		}
		val := string(kvEntry.Value())
		kv_map = append(kv_map, k+" = "+truncateText(val, goterm.Width()-15))
	}
	if err != nil {
		return []string{}, err
	}

	return kv_map, nil
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
		println(n)
		names = append(names, n)
	}

	return names, nil
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
	var lastCurPos int
	for {
		Clear()
		menu := cli.NewMenu(fmt.Sprintf("%v available buckets", server))

		for _, b := range buckets {
			menu.AddItem(b, b)
		}
		menu.AddItem("🚪 Quit", "quit")

		menu.CursorPos = lastCurPos

		choice := menu.Display()

		if choice == "quit" {
			break
		} else {
			// list all keys in a bucket
			lastCurPos = menu.CursorPos
			Clear()
			keys, err := GetAllKV(choice)
			if err != nil {
				panic(err)
			}
			t := goterm.Color(fmt.Sprintf("Content of %v\n", choice), goterm.YELLOW)
			fmt.Println(t)
			for _, k := range keys {
				println("  " + k)
			}
			fmt.Println("==================")
			fmt.Printf("%v entries", len(keys))
			cli.GetInput()
			Clear()
		}

	}

}
