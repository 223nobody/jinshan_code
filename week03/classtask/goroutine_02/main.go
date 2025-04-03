package main

import (
	"fmt"
	"math/rand"
)

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- rand.Intn(100) + 1
		}
		close(ch)
	}()
	for {
		if data, ok := <-ch; ok {
			fmt.Println(data)
		} else {
			break
		}
	}
}
