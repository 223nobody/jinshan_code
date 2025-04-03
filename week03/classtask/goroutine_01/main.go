package main

import (
	"fmt"
	"sync"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	totalSum := Square(nums, 3)
	fmt.Printf("分部计算的总和为: %d\n", totalSum)
}

func Square(nums []int, parts int) int {
	var wg sync.WaitGroup
	ch := make(chan int, parts)
	PartialSize := (len(nums) + parts - 1) / parts
	for i := 0; i < parts; i++ {
		start := i * PartialSize
		end := start + PartialSize
		if end > len(nums) {
			end = len(nums)
		}
		if start >= end {
			break
		}
		wg.Add(1)
		go func(slice []int) {
			defer wg.Done()
			sum := 0
			for _, num := range slice {
				sum += num * num
			}
			ch <- sum
		}(nums[start:end])
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	total := 0
	for partialsum := range ch {
		total += partialsum
	}
	return total
}
