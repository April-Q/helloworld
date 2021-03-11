package main

import "fmt"

func main() {
	fmt.Println("Result of Fib 9:", Fib(9))
	fmt.Println("Result of Fib 10:", Fib(10))
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}
