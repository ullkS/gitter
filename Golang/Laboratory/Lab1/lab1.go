package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
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

	var M int
	fmt.Print("Введите число потоков M: ")
	_, err = fmt.Scanln(&M)
	if err != nil {
		fmt.Println("Ошибка при чтении числа M:", err)
		return
	}

	// Генерация файла с числами от 1 до N
	err = generateFile(N)
	if err != nil {
		fmt.Println("Ошибка при генерации файла:", err)
		return
	}

	// Читаем числа из файла и перемещаем их в массив
	numbers, err := readNumbersFromFile("numbers.txt")
	if err != nil {
		fmt.Println("Ошибка при чтении чисел из файла:", err)
		return
	}
	// Последовательная обработка массива 1 поток
	start := time.Now()
	sequentialProcessing(numbers)
	sequentialProcessing := time.Since(start)
	
	// Многопоточная обработка массива
	start = time.Now()
	parallelProcessing(numbers, M)
	parallelProcessing := time.Since(start)
	
	fmt.Println("Последовательная обработка заняла:", sequentialProcessing)
	fmt.Println("Многопоточная обработка заняла:", parallelProcessing)
}

func generateFile(N int) error {
	file, err := os.Create("numbers.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 1; i <= N; i++ {
		_, err = file.WriteString(strconv.Itoa(i) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func readNumbersFromFile(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var numbers []int
	var num int
	for {
		_, err := fmt.Fscanf(file, "%d\n", &num)
		if err != nil {
			break
		}
		numbers = append(numbers, num)
	}

	return numbers, nil
}

func sequentialProcessing(numbers []int) {
	for _, num := range numbers {
		math.Pow(float64(num), 2)
	}
}

func parallelProcessing(numbers []int, M int) {
	var wg sync.WaitGroup
	chunkSize := int(math.Ceil(float64(len(numbers)) / float64(M)))

	for i := 0; i < M; i++ {
		start := i * chunkSize
		end := start + chunkSize

		wg.Add(1)
		go func(nums []int) {
			defer wg.Done()
			for _, num := range nums {
				math.Pow(float64(num), 2)
			}
		}(numbers[start:end])
	}

	wg.Wait()
}