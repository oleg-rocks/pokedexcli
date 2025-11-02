package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedox > ")
		scanner.Scan()
		input := scanner.Text()
		cleaned := cleanInput(input)
		fmt.Printf("Your command was: %s\n", cleaned[0])
	}
}

func cleanInput(text string) []string {
	lowText := strings.ToLower(text)
	words := strings.Fields(lowText)
	return words
}
