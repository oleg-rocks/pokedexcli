package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	lowText := strings.ToLower(text)
	words := strings.Fields(lowText)
	return words
}
