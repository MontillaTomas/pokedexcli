package main

import (
	"fmt"
	"strings"
)

func main() {
	var text string = "Hello, World!"
	fmt.Println(text)
	cleanedWords := cleanInput(text)
	fmt.Println("Returned words:", cleanedWords)
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}
