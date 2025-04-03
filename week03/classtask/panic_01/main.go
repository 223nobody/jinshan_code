package main

import "fmt"

func main() {
	var a, b int
	fmt.Scanln(&a, &b)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from: %v\n", r)
		}
	}()
	fmt.Println(divide(a, b))
}

func divide(a, b int) (result float64) {
	if b == 0 {
		panic("发生异常,被除数为0")
	}
	return float64(a) / float64(b)
}
