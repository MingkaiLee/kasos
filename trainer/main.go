package main

import "fmt"

type A struct {
	B int
}

func a() (p *A) {
	p.B = 1
	return
}

func main() {
	fmt.Println(a())
}
