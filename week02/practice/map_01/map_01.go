package main

import "fmt"

func main() {
	scores := map[string]int{
		"小明": 60,
		"小王": 70,
		"张三": 95,
		"李四": 98,
		"王五": 100,
		"张伟": 88,
	}
	for key, value := range scores {
		fmt.Printf("%s的成绩是 %d\n", key, value)

	}
}
