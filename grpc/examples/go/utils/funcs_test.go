package utils

import (
	"fmt"
	"testing"
)

func HelloWorld() {
	fmt.Println("Hello, World!")
}

func TestFunc(t *testing.T) {
	f, err := NewFunc(At[any])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f.Name())
	fmt.Println(f.FileLine())
	fmt.Println(f)
}
