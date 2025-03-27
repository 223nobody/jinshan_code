package main

import (
	"fmt"
)

func main() {
	arr := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice := arr[2:7]
	slice = append(slice, 11, 12, 13)
	slice = append(slice[:4], slice[5:]...)
	for i := range slice {
		slice[i] *= 2
	}
	fmt.Println(slice)
	fmt.Println(len(slice))
	fmt.Println(cap(slice))

}
