package main

import (
	"fmt"
)

func countCharacters(s string) rune {
	counter := make(map[rune]int)
	var a rune
	n := 0
	for _, char := range s {
		counter[char]++
	}
	for key, value := range counter {
		if value >= n {
			n = value
			a = key
		}
	}
	return a
}
func main() {
	var str string
	fmt.Scanln(&str)
	char := countCharacters(str)
	fmt.Printf("出现最多的字符为： %c", char)
}
