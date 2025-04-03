package main

import (
	"fmt"
)

func main() {
	var n int
	fmt.Scanln(&n)
	num := []int{1, 2, 3, 4, 5}
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from: %v\n", r)
		}
	}()
	fmt.Printf("索引位置为 %d 处的数组元素是 %d \n", n, accessArray(n, num))
}

func accessArray(n int, arr []int) int {
	if n >= len(arr) {
		panic("数组下标越界")
	}
	return arr[n]
}
