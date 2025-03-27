package main

import (
	"fmt"
)

func main() {

	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4, 5, 6}
	var combinedSlice, uniqueSlice []int
	combinedSlice = append(slice1, slice2...)
	for i := 0; i < len(combinedSlice); i++ {
		flag := false
		for j := i + 1; j < len(combinedSlice); j++ {
			if combinedSlice[i] == combinedSlice[j] {
				flag = true
				break
			}
		}
		if !flag {
			uniqueSlice = append(uniqueSlice, combinedSlice[i])
		}
	}
	fmt.Println(uniqueSlice)
}
