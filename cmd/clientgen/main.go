package main

import (
	"flag"
	"github.com/nsatyasrikar/client"
	"log"
)

func main() {
	data := flag.String("data", "", "generated JSON root")
	output := flag.String("output", ".", "client package root")
	check := flag.Bool("check", false, "check freshness")
	flag.Parse()
	if *data == "" {
		log.Fatal("-data is required")
	}
	if err := client.GenerateBindings(*data, *output, *check); err != nil {
		log.Fatal(err)
	}
}
