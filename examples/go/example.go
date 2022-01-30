package main

import "fmt"

func main() {
	a := 42

	for i := 0; i < 42; i++ {
		a += i
	}

	fmt.Println("Hello World", a)
	fmt.Println("sum of 1 + 2 is", sum(1, 2))
}

func sum(a, b int) int {
	return a + b
}
