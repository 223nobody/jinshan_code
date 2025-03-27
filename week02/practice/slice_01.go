package main

import "fmt"

func slice_02() {
	arr := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice := arr[2:7]
	slice = append(slice, 100)
	sum := 0
	for _, value := range slice {
		sum += value
	}
	fmt.Println(slice)
	fmt.Println(sum)

}
