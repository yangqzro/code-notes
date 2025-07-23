package main

import (
	"fmt"
	"goexamples/utils"
	"time"

	"google.golang.org/grpc/metadata"
)

func main() {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.DateTime),
		"grpc-timestamp", "example",
	)

	md.Set("grpc-timestamp", "example2") // This will panic if "grpc-timestamp" already exists
	fmt.Println(utils.String(md))
}
