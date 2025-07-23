package utils

import (
	"fmt"
	"testing"
)

type Hobby struct {
	Name string
}

type Person struct {
	Name    string
	name    string
	Age     int
	Hobby   []Hobby
	Test    []string
	TestMap map[string]string
}

func TestString(t *testing.T) {
	p1 := &Person{
		Name:    "1",
		Age:     1,
		Hobby:   []Hobby{{Name: "1"}, {Name: "2"}},
		Test:    []string{"1", "2"},
		TestMap: map[string]string{"1": "1", "2": "2"},
	}

	s := "test"
	sp := &s
	fmt.Println(String(&sp))
	fmt.Println(String(rune('我')))
	fmt.Println(String(byte('a')))
	fmt.Println(String(97))
	fmt.Println(String(GetPointElem))
	fmt.Println(String(p1))
	fmt.Println(String(map[byte]*Person{'a': p1, 'b': p1}))
	fmt.Println(String(make(chan int)))
}
