package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func ReshetoEratosthenes(n int) []int { // Вычеркивание по мод-му методу p^2
	primes := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		primes[i] = true
	}

	for p := 2; p*p <= n; p++ {
		if primes[p] {
			for i := 2 * 2; i <= n; i += p {
				primes[i] = false
			}
		}
	}

	result := []int{}
	for i := 2; i <= n; i++ {
		if primes[i] {
			result = append(result, i)
		}
	}

	return result
}

func ModifiledAlgorithmCheck(n int) bool {
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n <= 1 {
		return false
	}

	for i := 5; i <= int(math.Sqrt(float64(n))); i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}
func ModifiledAlgorithm(n int) []int { //Првоерка модифи-ым методом
	primes := []int{}
	for i := 2; i <= n; i++ {
		if ModifiledAlgorithmCheck(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

func DecompositionDataCheck(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := 5; i <= int(math.Sqrt(float64(n))); i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}

	return true
}
func DecompositionData(start, end int, wg *sync.WaitGroup, primes chan<- int) {
	defer wg.Done()

	for i := start; i <= end; i++ {
		if DecompositionDataCheck(i) {
			primes <- i
		}
	}
}
func StartDecomposition(numWorkers int, n int, primes chan int) { // №1 проверка многопотокм в декомпозиции
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		start := i*(n/numWorkers) + 1
		end := (i + 1) * (n / numWorkers)
		if i == numWorkers-1 {
			end = n
		}

		go DecompositionData(start, end, &wg, primes)
	}

	go func() {
		wg.Wait()
		close(primes)
	}()
	result := []int{}
	for prime := range primes {
		result = append(result, prime)
	}
}

func StartDecompositionNumeric(numWorkers int, n int, primes1 chan int) {
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := i*(n/numWorkers) + 1
		end := (i + 1) * (n / numWorkers)
		if i == numWorkers-1 {
			end = n
		}

		go DicompositionbasikNum(start, end, &wg, primes1)
	}

	go func() {
		wg.Wait()
		close(primes1)
	}()

	result := []int{}
	for prime := range primes1 {
		result = append(result, prime)
	}
}
func DecompositionNum(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := 5; i <= int(math.Sqrt(float64(n))); i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}

	return true
}
func DicompositionbasikNum(start, end int, wg *sync.WaitGroup, primes1 chan<- int) { // №2 декомпозиция набора простых чисел
	defer wg.Done()

	for i := start; i <= end; i++ {
		if DecompositionNum(i) {
			primes1 <- i
		}
	}
}

func StartPullPotok(N int, M int) { // №3 ПРименение пула потоков на старте
	// Создаем пул потоков
	var wg sync.WaitGroup
	jobs := make(chan int, N)
	results := make(chan int, N)

	// Запускаем горутины для обработки рабочих элементов
	for i := 0; i < M; i++ {
		wg.Add(1)
		go WorkerPullPotok(jobs, results, &wg)
	}

	// Передаем рабочие элементы в пул потоков
	for i := 2; i <= N; i++ {
		jobs <- i
	}

	close(jobs)
	wg.Wait()
	close(results)
	/*for prime := range results {
	fmt.Printf("%d is prime\n", prime) */

}
func WorkerPullPotok(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for number := range jobs {
		if CheckNumProstoe(number) {
			results <- number
		}
	}
}
func CheckNumProstoe(number int) bool {
	if number < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(number))); i++ {
		if number%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	N := 1000000     // Размер выборки
	NumWorkers := 12 // Количество потоков
	primes := make(chan int)
	primes1 := make(chan int)

	startTime := time.Now()
	ReshetoEratosthenes(N) // Можно присвоить и вывести массив.
	endTime := time.Now()

	startTime2 := time.Now()
	ModifiledAlgorithm(N)
	endTime2 := time.Now()

	startTime3 := time.Now()
	StartDecomposition(NumWorkers, N, primes)
	endTime3 := time.Now()

	startTime4 := time.Now()
	StartDecompositionNumeric(NumWorkers, N, primes1)
	endTime4 := time.Now()

	startTime5 := time.Now()
	StartPullPotok(N, NumWorkers)
	endTime5 := time.Now()

	fmt.Println("Время выполнения(Решето Эратосфена):", endTime.Sub(startTime))
	fmt.Println("Время выполнения(Модифицированный последовательный алгоритм поиска):", endTime2.Sub(startTime2))
	fmt.Println("Время выполнения(ПА декомпозиция по данным):", endTime3.Sub(startTime3))
	fmt.Println("Время выполнения(ПА декомпозиция по простым числам):", endTime4.Sub(startTime4))
	fmt.Println("Время выполнения(ПА применение пула потоков):", endTime5.Sub(startTime5))
}
