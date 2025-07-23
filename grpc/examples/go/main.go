package main

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	fmt.Printf("runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{\"Greeter\", \"SayHello\"}, \"\")): %v\n", runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"Greeter", "SayHello"}, "")))
}
