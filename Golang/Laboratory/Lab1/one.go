package main

import (
	"fmt"
	"math"
	"time"
)

func main() {

	var N int
	fmt.Print("Введите число N: ")
	_, err := fmt.Scanln(&N)
	if err != nil {
		fmt.Println("Ошибка при чтении числа N:", err)
		return
	}
	numbers := []int{}

	for i := 1; i < N; i++ {
		numbers = append(numbers, i)
	}

	start := time.Now()
	sequentialProcessingg(numbers)
	sequentialProcessing := time.Since(start)

	fmt.Println("Последовательная обработка заняла:", sequentialProcessing)
}

func sequentialProcessingg(numbers []int) {
	for _, num := range numbers {
		math.Pow(float64(num), 2)
	}
}
