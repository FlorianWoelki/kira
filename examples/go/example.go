package main

import "fmt"

func main() {
	a := 42

	for i := 0; i < 42; i++ {
		a += i
	}

	fmt.Println("Hello World", a)
}
