package main

import (
	"fmt"
	"os"
)

func f(i int) (int, error) {
	os.Exit(1)
	return i * 2, nil
}

func main() {
	os.Exit(1) // want "direct use of os.Exit in main function is prohibited"
	res, _ := f(5)
	fmt.Println(res)

	defer os.Exit(1) // want "direct use of os.Exit in main function is prohibited"
	go os.Exit(1)    // want "direct use of os.Exit in main function is prohibited"
}
