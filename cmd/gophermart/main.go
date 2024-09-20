package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("=========")

	panic(errors.New("WTF"))
}
