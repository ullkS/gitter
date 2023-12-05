package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func main() {

	var N int = 100000000

	M := [8]int{2, 3, 4, 5, 10, 20, 30, 100}
	numbers := []int{}

	for i := 1; i < N; i++ {
		numbers = append(numbers, i)
	}
	for _, m := range M {
		start := time.Now()
		parallelProcessingg(numbers, m)
		parallelProcessing := time.Since(start)
		fmt.Println("Многопоточная обработка заняла:", m, parallelProcessing)
	}
}

func parallelProcessingg(numbers []int, M int) {
	var wg sync.WaitGroup
	chunkSize := int(math.Ceil(float64(len(numbers)) / float64(M)))

	for i := 0; i < M; i++ {
		start := i * chunkSize
		end := start + chunkSize

		wg.Add(1)
		go func(nums []int) {
			defer wg.Done()
			for count, num := range nums {
				num = num * count
			}
		}(numbers[start:end])
	}
	wg.Wait()
}
