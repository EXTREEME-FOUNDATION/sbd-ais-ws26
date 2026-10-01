// Create a hello world program that prints the contents of a struct using a method with a receiver type
package main

import "fmt"

// Strukt
type Strukt struct {
	something string
	Test123   string
}

func (s Strukt) Print() { // reciever type is Strukt
	fmt.Println(s.something)
	fmt.Println(s.Test123)
}

func main() {
	s := Strukt{something: "Hello World", Test123: "This is a test"}
	s.Print()
}
