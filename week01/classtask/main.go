package main

import "fmt"

func Judge(num int) bool {
	if num < 0 {
		return false
	}

	if num >= 0 && num < 10 {
		return true
	}
	reversed := 0
	original := num
	for num > 0 {
		digit := num % 10
		reversed = reversed*10 + digit
		num = num / 10
	}

	return original == reversed
}

func main() {
	var num int
	fmt.Print("请输入一个整数：")
	fmt.Scan(&num)

	if Judge(num) {
		fmt.Printf("%d 是回文数。\n", num)
	} else {
		fmt.Printf("%d 不是回文数。\n", num)
	}
}
