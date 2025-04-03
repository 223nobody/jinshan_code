package main

import "fmt"

func doubleValues(num *[]int) {
	for p := range *num {
		(*num)[p] *= 2
	}
}

func main() {
	arr := []int{1, 2, 3, 4, 5}
	doubleValues(&arr)
	fmt.Println(arr)
}
